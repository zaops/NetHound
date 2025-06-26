from tkinter import INSERT
import duckdb

# with duckdb.connect("file.db") as con:
#     con.sql("CREATE TABLE test (i INTEGER)")
#     con.sql("INSERT INTO test VALUES (42)")
#     con.table("test").show()
#     # the context manager closes the connection automatically

con = duckdb.connect("file.db")

con.sql("INSERT INTO temp_phone_numbers VALUES ('136410093333')")

con.table("temp_phone_numbers").show()

con.close()