DROP TABLE IF EXISTS `face_record`;
CREATE TABLE `face_record` (
    `id` INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT ,
    `device_id` VARCHAR(40) NOT NULL ,
    `device_latitude` FLOAT NOT NULL ,
    `device_longtitude` FLOAT NOT NULL ,
    `img_path` VARCHAR(255) NOT NULL
);

DROP TABLE IF EXISTS `car_record`;
CREATE TABLE `car_record` (
    `id` INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT ,
    `device_id` VARCHAR(40) NOT NULL ,
    `device_latitude` FLOAT NOT NULL ,
    `device_longtitude` FLOAT NOT NULL ,
    `car_number` VARCHAR(255) NOT NULL ,
    `img_path` VARCHAR(255) NOT NULL
);