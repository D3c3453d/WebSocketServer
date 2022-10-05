CREATE TABLE IF NOT EXISTS users
(
    id            SERIAL       NOT NULL,
    username      VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255),

    PRIMARY KEY (id)
) WITHOUT OIDS;

CREATE TABLE IF NOT EXISTS friends
(
    id          SERIAL       NOT NULL,
    user_id     INT          NOT NULL,
    friend_id   INT          NOT NULL,

    PRIMARY KEY (id),
    CONSTRAINT friends_users_user_id_foreign FOREIGN KEY (user_id) REFERENCES users (id),
    CONSTRAINT friends_users_friend_id_foreign FOREIGN KEY (friend_id) REFERENCES users (id)
) WITHOUT OIDS;

