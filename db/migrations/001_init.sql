-- create db tables for the users (viewer (by default) , analyst and admin) and table for records  

CREATE TABLE IF NOT EXISTS users (
    id  TEXT PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    role TEXT NOT NULL DEFAULT 'viewer',
    is_active INTEGER NOT NULL DEFAULT 1,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);


-- table for finance records 

CREATE TABLE IF NOT EXISTS financial_records (
      id TEXT PRIMARY KEY,
      user_id TEXT NOT NULL,
      amount REAL NOT NULL,
      type TEXT NOT NULL CHECK(type IN('income' , "expense")),
      category TEXT NOT NULL,
      date TEXT NOT NULL,
      notes TEXT,
      created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
      updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
      FOREIGN KEY (user_id) REFERENCES users(id)
);



 -- some index on table for fast lookups  

 CREATE INDEX IF NOT EXISTS idx_records_user_id ON financial_records(user_id);

 CREATE INDEX IF NOT EXISTS idx_records_date ON financial_records(date);

 CREATE INDEX IF NOT EXISTS idx_records_type ON financial_records(type);


 CREATE INDEX IF NOT EXISTS idx_records_category ON financial_records(category);








