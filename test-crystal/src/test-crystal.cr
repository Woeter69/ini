require "db"
require "sqlite3"

# test-crystal DB App
DB.open "sqlite3:./test.db" do |db|
  db.exec "CREATE TABLE IF NOT EXISTS users (name TEXT)"
  db.exec "INSERT INTO users VALUES (?)", "User from test-crystal"
  
  puts "Connected and inserted a record into SQLite for test-crystal"
end
