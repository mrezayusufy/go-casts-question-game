CREATE TABLE categories {
  id          int PRIMARY KEY AUTO_INCREMENT NOT NULL,
  name        varchar(255) NOT NULL,
  description varchar(255),
  created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
}