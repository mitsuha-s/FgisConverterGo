CREATE TABLE config (
    signcipher varchar(10) CHECK(length(signcipher) > 0),
    ambientmode varchar(10) DEFAULT 'empty' CHECK(ambientmode IN('empty', 'file', 'random'))
);

CREATE TABLE methods (
    id integer PRIMARY KEY,
    name varchar(20)
);

CREATE TABLE mitypes (
    regnum varchar(8) UNIQUE NOT NULL CHECK (length(regnum) == 8),
    method integer NOT NULL,
    FOREIGN KEY (method) REFERENCES methods
);

CREATE TABLE testdevicemi (
    id integer PRIMARY KEY,
    regnum varchar(8) NOT NULL CHECK (length(regnum) == 8),
    number varchar(10) UNIQUE NOT NULL CHECK (length(number) > 0)
);

CREATE TABLE testdevices (
    regnum varchar(20) UNIQUE NOT NULL CHECK (length(regnum) > 0),
    stopwatch integer,
    thermometer integer,
    thermohygrometer integer,
    FOREIGN KEY (stopwatch) REFERENCES testdevicemi,
    FOREIGN KEY (thermometer) REFERENCES testdevicemi,
    FOREIGN KEY (thermohygrometer) REFERENCES testdevicemi,
    CONSTRAINT mi CHECK(stopwatch != thermohygrometer AND stopwatch != thermometer AND thermometer != thermohygrometer)
);