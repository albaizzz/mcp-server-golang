package tools
package main

import (
    "bufio"
    "context"
    "database/sql"
    "encoding/csv"
    "flag"
    "fmt"
    "log"
    "os"
    "time"

    _ "github.com/microsoft/go-mssqldb"
)

func main() {
    csvPath := flag.String("csv", "kodepos.csv", "path to csv")
    dsn := flag.String("dsn", "sqlserver://sa:YourStrong!Passw0rd@localhost:1433?database=mcp_kodepos&encrypt=disable", "sqlserver dsn")
    batch := flag.Int("batch", 1000, "batch size")
    flag.Parse()

    db, err := sql.Open("sqlserver", *dsn)
    if err != nil { log.Fatal(err) }
    defer db.Close()

    f, err := os.Open(*csvPath)
    if err != nil { log.Fatal(err) }
    defer f.Close()

    r := csv.NewReader(bufio.NewReader(f))
    r.ReuseRecord = true
    r.FieldsPerRecord = -1

    // skip header if present: detect by first record's first field == "kodepos"
    first, err := r.Read()
    if err != nil { log.Fatal(err) }
    var records [][]string
    if len(first) > 0 && first[0] != "kodepos" {
        records = append(records, first)
    }

    for {
        rec, err := r.Read()
        if err != nil {
            if err.Error() == "EOF" { break }
            log.Fatal(err)
        }
        records = append(records, rec)
        if len(records) >= *batch {
            if err := insertBatch(db, records); err != nil { log.Fatal(err) }
            records = records[:0]
        }
    }
    if len(records) > 0 {
        if err := insertBatch(db, records); err != nil { log.Fatal(err) }
    }
    fmt.Println("Import done")
}

func insertBatch(db *sql.DB, recs [][]string) error {
    ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
    defer cancel()

    tx, err := db.BeginTx(ctx, nil)
    if err != nil { return err }

    stmt, err := tx.PrepareContext(ctx, `
        INSERT INTO dbo.kodepos (kodepos, kelurahan, kecamatan, kota_kab, provinsi, latitude, longitude)
        VALUES (@p1, @p2, @p3, @p4, @p5, @p6, @p7)
    `)
    if err != nil { _ = tx.Rollback(); return err }
    defer stmt.Close()

    for _, r := range recs {
        var p [7]any
        for i := 0; i < 7 && i < len(r); i++ {
            if r[i] == "" && (i == 5 || i == 6) {
                p[i] = nil
            } else {
                p[i] = r[i]
            }
        }
        if _, err := stmt.ExecContext(ctx, p[0], p[1], p[2], p[3], p[4], p[5], p[6]); err != nil {
            _ = tx.Rollback(); return err
        }
    }
    return tx.Commit()
}
