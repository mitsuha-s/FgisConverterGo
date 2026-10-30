CREATE TABLE config (
    signcipher varchar(10),
    ambientmod varchar(10) CHECK(ambientmod IN('empty', 'file', 'random'))
);

CREATE TABLE methods (
    id integer PRIMARY KEY,
    name varchar(20)
);

CREATE TABLE mitypes (
    regnum varchar(8) UNIQUE,
    method integer,
    FOREIGN KEY (method) REFERENCES methods
);

CREATE TABLE testdevicemi (
    id integer PRIMARY KEY,
    regnum varchar(8),
    number varchar(10)
);

CREATE TABLE testdevices (
    regnum varchar(20),
    stopwatch integer,
    thermometer integer,
    thermohygrometer integer,
    FOREIGN KEY (stopwatch) REFERENCES testdevicemi,
    FOREIGN KEY (thermometer) REFERENCES testdevicemi,
    FOREIGN KEY (thermohygrometer) REFERENCES testdevicemi,
    CONSTRAINT mi CHECK(stopwatch != thermohygrometer AND stopwatch != thermometer AND thermometer != thermohygrometer)
);