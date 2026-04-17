CREATE TABLE users (
  id                INT PRIMARY KEY NOT NULL AUTO_INCREMENT,
  name              VARCHAR(255) NOT NULL DEFAULT "reza",
  password          VARCHAR(255) NOT NULL,
  phone_number      VARCHAR(255) NOT NULL DEFAULT "090" UNIQUE,
  created_at        TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
)
INSERT INTO users(id, name, phone_number) values(1, "reza", "09030072667")