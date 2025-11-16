CREATE TABLE transactions (
    id INT PRIMARY KEY AUTO_INCREMENT,
    user_public_id VARCHAR(100) NOT NULL,
    product_id INT NOT NULL,
    product_price INT NOT NULL,
    amount INT NOT NULL,
    sub_total INT NOT NULL,
    platform_fee INT DEFAULT 0,
    grand_total INT NOT NULL,
    status int NOT NULL,
    product_snapshot JSON,
    created_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    update_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);