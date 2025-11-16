package sql

import (
	"log"

	"github.com/jmoiron/sqlx"
)

func RunMigrations(db *sqlx.DB) {
	authTable := `
	CREATE TABLE IF NOT EXISTS auth (
		id INT AUTO_INCREMENT PRIMARY KEY,
		public_id VARCHAR(255) NOT NULL,
		email VARCHAR(255) NOT NULL UNIQUE,
		password VARCHAR(255) NOT NULL,
		role ENUM('super_admin', 'admin', 'user') NOT NULL DEFAULT 'user',
		created_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		update_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
	);
	`

	productTable := `
	CREATE TABLE IF NOT EXISTS products (
		id INT AUTO_INCREMENT PRIMARY KEY,
		sku VARCHAR(100) NOT NULL,
		name VARCHAR(100) NOT NULL,
		stock INT NOT NULL DEFAULT 0,
		price INT NOT NULL DEFAULT 0,
		created_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		update_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
	);
	`

	transactionTable := `
	CREATE TABLE IF NOT EXISTS transactions (
		id INT AUTO_INCREMENT PRIMARY KEY,
		user_public_id VARCHAR(100) NOT NULL,
		product_id INT NOT NULL,
		product_price INT NOT NULL,
		amount INT NOT NULL,
		sub_total INT NOT NULL,
		platform_fee INT DEFAULT 0,
		grand_total INT NOT NULL,
		status INT NOT NULL,
		product_snapshot JSON,
		created_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		update_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		FOREIGN KEY (product_id) REFERENCES products(id)
	);
	`

	queries := []string{authTable, productTable, transactionTable}

	for _, query := range queries {
		_, err := db.Exec(query)
		if err != nil {
			log.Fatalf("❌ Migration failed: %v", err)
		}
	}

	log.Println("✅ Database migration completed successfully.")
}