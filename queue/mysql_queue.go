package queue

import (
	"database/sql"
	"errors"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// MySQLQueue MySQL 队列实现
type MySQLQueue struct {
	db *sql.DB
}

// NewMySQLQueue 创建 MySQL 队列实例
func NewMySQLQueue(dsn string) (*MySQLQueue, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	// 测试连接
	if err := db.Ping(); err != nil {
		return nil, err
	}

	// 创建消息表
	if err := createMessageTable(db); err != nil {
		return nil, err
	}

	return &MySQLQueue{db: db}, nil
}

// createMessageTable 创建消息表
func createMessageTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS messages (
		id VARCHAR(36) PRIMARY KEY,
		type VARCHAR(50) NOT NULL,
		body TEXT NOT NULL,
		status TINYINT DEFAULT 0,
		create_at DATETIME NOT NULL,
		retry_count INT DEFAULT 0,
		next_retry_at DATETIME NOT NULL,
		INDEX idx_status_next_retry (status, next_retry_at)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
	`

	_, err := db.Exec(query)
	return err
}

// Produce 生产消息
func (q *MySQLQueue) Produce(message Message) error {
	query := `
	INSERT INTO messages (id, type, body, status, create_at, retry_count, next_retry_at)
	VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	_, err := q.db.Exec(
		query,
		message.ID,
		message.Type,
		message.Body,
		message.Status,
		message.CreateAt,
		message.RetryCount,
		message.NextRetryAt,
	)

	return err
}

// Consume 消费消息
func (q *MySQLQueue) Consume() (Message, error) {
	var message Message

	// 使用事务确保消息的原子性
	tx, err := q.db.Begin()
	if err != nil {
		return message, err
	}

	// 查找待处理或需要重试的消息
	query := `
	SELECT id, type, body, status, create_at, retry_count, next_retry_at
	FROM messages
	WHERE status IN (0, 2) AND next_retry_at <= ?
	ORDER BY next_retry_at ASC
	LIMIT 1
	FOR UPDATE SKIP LOCKED
	`

	err = tx.QueryRow(query, time.Now()).Scan(
		&message.ID,
		&message.Type,
		&message.Body,
		&message.Status,
		&message.CreateAt,
		&message.RetryCount,
		&message.NextRetryAt,
	)

	if err != nil {
		tx.Rollback()
		if errors.Is(err, sql.ErrNoRows) {
			return message, nil // 没有消息
		}
		return message, err
	}

	// 更新消息状态为处理中
	updateQuery := `
	UPDATE messages
	SET status = 1
	WHERE id = ?
	`

	_, err = tx.Exec(updateQuery, message.ID)
	if err != nil {
		tx.Rollback()
		return message, err
	}

	if err := tx.Commit(); err != nil {
		return message, err
	}

	return message, nil
}

// UpdateMessageStatus 更新消息状态
func (q *MySQLQueue) UpdateMessageStatus(id string, status int8, retryCount int, nextRetryAt time.Time) error {
	query := `
	UPDATE messages
	SET status = ?, retry_count = ?, next_retry_at = ?
	WHERE id = ?
	`

	_, err := q.db.Exec(query, status, retryCount, nextRetryAt, id)
	return err
}
