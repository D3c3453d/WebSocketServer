INSERT INTO users (username) VALUES ('aaa');
INSERT INTO users (username) VALUES ('bbb');
INSERT INTO users (username) VALUES ('ccc');

INSERT INTO friends (user_id, friend_id) VALUES (1, 2);
INSERT INTO friends (user_id, friend_id) VALUES (2, 1);
INSERT INTO friends (user_id, friend_id) VALUES (2, 3);
INSERT INTO friends (user_id, friend_id) VALUES (3, 2);