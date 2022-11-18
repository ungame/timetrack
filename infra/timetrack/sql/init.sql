CREATE TABLE IF NOT EXISTS categories (
    id INT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(50) NOT NULL UNIQUE,
    description TEXT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
)
ENGINE = INNODB
DEFAULT CHARSET = UTF8;

INSERT INTO categories (id, name, description) VALUES
(1, 'meeting', 'daily, 1:1, refinement, retro'),
(2, 'coding', 'features, bugs, tests'),
(3, 'review', 'pull requests'),
(4, 'pm', 'private messages'),
(5, 'task', '')
;

CREATE TABLE IF NOT EXISTS activities (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    category_id INT NOT NULL,
    description TEXT NOT NULL,
    status CHAR(1) DEFAULT '1',
    started_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    finished_at TIMESTAMP NULL,
    CONSTRAINT activities_category_id_fk
    FOREIGN KEY (category_id)
    REFERENCES categories(id)
)
ENGINE = INNODB
DEFAULT CHARSET = UTF8;