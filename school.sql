CREATE TABLE subjects (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    subject TEXT NOT NULL DEFAULT ""
);

CREATE TABLE teachers (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    fio TEXT NOT NULL DEFAULT "",
    phone VARCHAR(32) NOT NULL DEFAULT ""
);

CREATE TABLE teachers_subjects (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    teacher_id INTEGER NOT NULL,
    subject_id INTEGER NOT NULL
);

CREATE TABLE students (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    fio TEXT NOT NULL DEFAULT "",
    phone VARCHAR(32) NOT NULL DEFAULT "",
    birthday CHAR(8) NOT NULL DEFAULT "",
    adress TEXT NOT NULL DEFAULT "",
    class_id INTEGER NOT NULL
);

CREATE TABLE classes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    level INTEGER NOT NULL,
    letter CHAR(1) NOT NULL
);

CREATE TABLE teachers_classes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    teacher_id INTEGER NOT NULL,
    classe_id INTEGER NOT NULL
);

CREATE INDEX classe_teachers ON teachers_classes (classe_id);
CREATE INDEX teacher_classes ON teachers_classes (teacher_id);
CREATE INDEX teacher_subjects ON teachers_subjects (teacher_id); 
CREATE INDEX list_students_in_class ON students (class_id);