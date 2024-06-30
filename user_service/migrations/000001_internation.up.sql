CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS "branch" (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    address VARCHAR(255),
    phone VARCHAR(20) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at INTEGER DEFAULT 0
);

CREATE TABLE IF NOT EXISTS "group" (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    teacherId UUID,
    supportTeacherId UUID,
    branchId UUID,
    type VARCHAR(50) CHECK (type IN ('beginner', 'elementary', 'intermediate', 'ielts')),
    FOREIGN KEY (teacherId) REFERENCES teacher(id),
    FOREIGN KEY (supportTeacherId) REFERENCES support_teacher(id),
    FOREIGN KEY (branchId) REFERENCES branch(id),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at INTEGER DEFAULT 0
);

CREATE TABLE IF NOT EXISTS "journal" (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    fromDate DATE,
    toDate DATE,
    groupId UUID,
    FOREIGN KEY (groupId) REFERENCES group(id)
    studentsCount INT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at INTEGER DEFAULT 0
);

CREATE TABLE IF NOT EXISTS "support_teacher" (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    login VARCHAR(255) UNIQUE NOT NULL,
    fullname VARCHAR(255),
    phone VARCHAR(20) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    salary INTEGER,
    ieltsScore FLOAT,  
    ieltsAttemptCount INTEGER,
    branchId UUID,
    FOREIGN KEY (branchId) REFERENCES branch(id),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at INTEGER DEFAULT 0
);

CREATE TABLE IF NOT EXISTS "teacher" (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    login VARCHAR(255) UNIQUE NOT NULL,
    fullname VARCHAR(255),
    phone VARCHAR(20) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    salary INTEGER,
    ieltsScore FLOAT,  
    ieltsAttemptCount INTEGER,
    supportTeacherId UUID,
    branchId UUID,
    FOREIGN KEY (supportTeacherId) REFERENCES support_teacher(id),
    FOREIGN KEY (branchId) REFERENCES branch(id),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at INTEGER DEFAULT 0
);

CREATE TABLE IF NOT EXISTS "administration" (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    login VARCHAR(255) UNIQUE NOT NULL,
    fullname VARCHAR(255),
    phone VARCHAR(20) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    salary INTEGER,
    ieltsScore FLOAT,  
    branchId UUID,
    FOREIGN KEY (branchId) REFERENCES branch(id),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at INTEGER DEFAULT 0
);

CREATE TABLE IF NOT EXISTS "student" (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    login VARCHAR(255) UNIQUE NOT NULL,
    fullname VARCHAR(255),
    phone VARCHAR(20) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    groupName VARCHAR(255),
    branchId UUID,
    FOREIGN KEY (branchId) REFERENCES branch(id),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at INTEGER DEFAULT 0
);

CREATE TABLE IF NOT EXISTS "manager" (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    login VARCHAR(255) UNIQUE NOT NULL,
    fullname VARCHAR(255),
    salary INTEGER,
    phone VARCHAR(20) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    branchId UUID,
    FOREIGN KEY (branchId) REFERENCES branch(id),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at INTEGER DEFAULT 0
);


CREATE TABLE IF NOT EXISTS "schedule" (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    journalId UUID,
    date DATE,
    startTime TIME,
    endTime TIME,
    lesson VARCHAR(255),
    FOREIGN KEY (journalId) REFERENCES journal(id),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at INTEGER DEFAULT 0
);

CREATE TABLE IF NOT EXISTS "task" (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    scheduleId UUID,
    label VARCHAR(255),
    deadlineDate DATE,
    deadlineTime TIME,
    score INT,
    FOREIGN KEY (scheduleId) REFERENCES schedule(id),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at INTEGER DEFAULT 0
);

CREATE TABLE IF NOT EXISTS "student_task" (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    taskId UUID,
    studentId UUID,
    FOREIGN KEY (taskId) REFERENCES task(id),
    FOREIGN KEY (studentId) REFERENCES student(id),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at INTEGER DEFAULT 0
);

CREATE TABLE IF NOT EXISTS "event" (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    assignStudent VARCHAR(255),
    topic VARCHAR(255),
    startTime TIME,
    date DATE,
    branchId UUID,
    FOREIGN KEY (branchId) REFERENCES branch(id),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at INTEGER DEFAULT 0
);

CREATE TABLE IF NOT EXISTS "event_student" (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    eventId UUID,
    studentId UUID,
    FOREIGN KEY (eventId) REFERENCES event(id),
    FOREIGN KEY (studentId) REFERENCES student(id),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at INTEGER DEFAULT 0
);

CREATE TABLE IF NOT EXISTS "student_payment" (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    studentId UUID,
    groupId UUID,
    paidSum DECIMAL(10, 2),
    administrationId UUID,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at INTEGER DEFAULT 0,
    FOREIGN KEY (studentId) REFERENCES "student"(id),
    FOREIGN KEY (groupId) REFERENCES "group"(id),
    FOREIGN KEY (administrationId) REFERENCES "administration"(id)
);