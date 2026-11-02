CREATE TABLE config (
    signcipher varchar(10) CHECK(length(signcipher) > 0),
    ambientmode varchar(10) DEFAULT 'empty' CHECK(ambientmode IN('empty', 'file', 'random'))
);

CREATE TABLE methods (
    id integer PRIMARY KEY,
    name varchar(20) NOT NULL
);

CREATE TABLE mitypes (
    regnum varchar(8) UNIQUE NOT NULL CHECK (length(regnum) == 8),
    method integer NOT NULL,
    FOREIGN KEY (method) REFERENCES methods
);

CREATE TABLE testdevicemi (
    regnum varchar(8) NOT NULL CHECK (length(regnum) == 8),
    number varchar(10) NOT NULL CHECK (length(number) > 0),
    testdevice integer,
    FOREIGN KEY (testdevice) REFERENCES testdevices,
    PRIMARY KEY (regnum, number) 
);

CREATE TABLE testdevices (
    id integer PRIMARY KEY,
    regnum varchar(20) UNIQUE NOT NULL CHECK (length(regnum) > 0)
);