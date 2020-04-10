SET SQL_SAFE_UPDATES = 0;
UPDATE user_profile SET weight = 30 WHERE weight < 30;
UPDATE user_profile SET weight = 500 WHERE weight > 500;
UPDATE user_profile SET height = 50 WHERE height < 50;
UPDATE user_profile SET height = 250 WHERE height > 250;
UPDATE user_profile SET birthday = '1940-01-01 00:00:00' WHERE TIMESTAMPDIFF(YEAR,birthday,CURDATE()) > 120;
UPDATE user_profile SET birthday = '1940-01-01 00:00:00' WHERE birthday = '0000-00-00 00:00:00';
