-- Create the Ads table with min_age and max_age (default values)
CREATE TABLE Ads (
  id INT PRIMARY KEY AUTO_INCREMENT,
  title VARCHAR(255) NOT NULL,
  start_at DATETIME NOT NULL,
  end_at DATETIME NOT NULL,
  min_age INT DEFAULT 1,  -- Default minimum age (1)
  max_age INT DEFAULT 100 -- Default maximum age (100)
--   FOREIGN KEY (id) REFERENCES AdConditions(ad_id)  -- Optional foreign key for future reference
);

-- Create the Conditions table for non-age targeting options
CREATE TABLE Conditions (
  id INT PRIMARY KEY AUTO_INCREMENT,
  type ENUM('gender', 'country', 'platform')  NOT NULL,
  value VARCHAR(255) DEFAULT NULL
);

-- Create the AdConditions table for linking ads and conditions
CREATE TABLE AdConditions (
  id INT PRIMARY KEY AUTO_INCREMENT,
  ad_id INT NOT NULL,
  condition_id INT NOT NULL,
  FOREIGN KEY (ad_id) REFERENCES Ads(id),
  FOREIGN KEY (condition_id) REFERENCES Conditions(id)
);

-- Add indexes for efficient querying (recommended)
ALTER TABLE Ads ADD INDEX (start_at, end_at);
ALTER TABLE Conditions ADD INDEX (type);

-- Insert gender conditions
INSERT INTO Conditions (type, value) VALUES ('gender', 'M');
INSERT INTO Conditions (type, value) VALUES ('gender', 'F');

-- Insert country conditions
INSERT INTO Conditions (type, value) VALUES ('country', 'TW');
INSERT INTO Conditions (type, value) VALUES ('country', 'JP');

-- Insert platform conditions
INSERT INTO Conditions (type, value) VALUES ('platform', 'android');
INSERT INTO Conditions (type, value) VALUES ('platform', 'ios');
INSERT INTO Conditions (type, value) VALUES ('platform', 'web');
