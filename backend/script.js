const sqlite3 = require('sqlite3').verbose();
const fs = require('fs');
const path = require('path');

// Use an environment variable for the database path, with a default fallback
const dbPath = process.env.SQLITE_DB_PATH || path.join(__dirname, 'database.sqlite');

console.log('Attempting to open database at:', dbPath);

// Check if the file exists
if (!fs.existsSync(dbPath)) {
    console.error('Database file does not exist at the specified path:', dbPath);
    process.exit(1);
}

const db = new sqlite3.Database(dbPath, sqlite3.OPEN_READONLY, (err) => {
  if (err) {
    console.error('Error opening database:', err.message);
    console.error('Database path:', dbPath);
    console.error('Current working directory:', process.cwd());
    process.exit(1);
  }
  console.log('Connected to the SQLite database.');
});

function getTableNames() {
  return new Promise((resolve, reject) => {
    const query = "SELECT name FROM sqlite_master WHERE type='table' AND name NOT LIKE 'sqlite_%'";
    db.all(query, [], (err, rows) => {
      if (err) {
        reject(err);
      } else {
        resolve(rows.map(row => row.name));
      }
    });
  });
}

function getTableData(tableName) {
  return new Promise((resolve, reject) => {
    db.all(`SELECT * FROM ${tableName}`, [], (err, rows) => {
      if (err) {
        reject(err);
      } else {
        resolve({
          tableName,
          columns: rows.length > 0 ? Object.keys(rows[0]) : [],
          rows: rows
        });
      }
    });
  });
}

async function exportDatabase() {
  try {
    const tables = await getTableNames();
    const databaseContent = {};

    for (const table of tables) {
      const tableData = await getTableData(table);
      databaseContent[table] = tableData;
    }

    fs.writeFileSync('real-db.json', JSON.stringify(databaseContent, null, 2));
    console.log('Database content exported to real-db.json successfully!');
  } catch (error) {
    console.error('Error exporting database:', error);
  } finally {
    db.close((err) => {
      if (err) {
        console.error('Error closing database:', err.message);
      } else {
        console.log('Closed the database connection.');
      }
    });
  }
}

exportDatabase();