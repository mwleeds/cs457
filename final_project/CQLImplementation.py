#!/usr/bin/python3

# This script implements a subset of CQL queries (which also happen to be SQL queries) 
# for use on the A.txt, B.txt, and C.txt files.
# Acceptable queries look like: "SELECT ... FROM ... JOIN ... WHERE ..."

import pprint
import re
from prettytable import PrettyTable

class CQLImplementation(object):
    def __init__(self):
        self.tables = {}
        counter = 0
        for filename in ('A.txt', 'B.txt', 'C.txt'):
            tablename = filename.split('.')[0]
            with open(filename) as f:
                lines = f.readlines()
                columns = [tablename + str(i+1) for i in range(len(lines[0].split()))] + ['writetime']
                rows = []
                for line in lines:
                    counter += 1
                    rows.append([int(s) for s in line.split()] + [counter])
                self.tables[tablename] = {'columns': columns, 'rows': rows}

    def print_db(self):
        pprint.pprint(self.tables)

    def compare(self, val1, op, val2):
        if op == '<': return val1 < val2
        elif op == '>': return val1 > val2
        elif op == '<=': return val1 <= val2
        elif op == '>=': return val1 >= val2
        elif op == '=': return val1 == val2
        elif op == '<>': return val1 != val2
        else: raise Exception('Invalid comparison operator passed to compare')

    def execute(self, query):
        # this regex matches a SELECT query on tables A, B, or C with optional JOIN and WHERE clauses
        select_re = r'^SELECT\s+(?P<cols>.+?)\s+' + \
                    r'FROM\s+(?P<tables>[A-C])' + \
                    r'((\s+JOIN\s+(?P<join_table>.+?)\s+ON\s+(?P<join_condition>.+?))?\s+' + \
                    r'WHERE\s+(?P<where_clause>.+?))?;$'
        tokens = re.match(select_re, query)
        if tokens == None:
            raise Exception('Invalid query entered')
        # get tokens
        cols = tokens.group('cols')
        tables = tokens.group('tables')
        join_table = tokens.group('join_table')
        join_condition = tokens.group('join_condition')
        where_clause = tokens.group('where_clause')
        if where_clause:
            where_re = r'(?P<col1>.+?)(?P<op1><|>|<=|>=|=|<>)(?P<val1>\d+?)' + \
                       r'(\s+(?P<binop>AND|OR)\s+' + \
                       r'(?P<col2>.+?)(?P<op2><|>|<=|>=|=|<>)(?P<val2>\d+?))?$'
            where_tokens = re.match(where_re, where_clause)
            if where_tokens == None:
                raise Exception('Invalid WHERE clause')
            col1, op1, val1 = where_tokens.group('col1'), where_tokens.group('op1'), where_tokens.group('val1')
            col2, op2, val2 = where_tokens.group('col2'), where_tokens.group('op2'), where_tokens.group('val2')
            binop = where_tokens.group('binop')
        table_columns = self.tables[tables]['columns']
        if cols == '*': cols = table_columns
        else: cols = cols.replace(' ', '').split(',')
        t = PrettyTable(cols)
        rows_returned = 0
        # execute query
        for row in self.tables[tables]['rows']:
            #TODO implement JOIN
            include = True
            if where_clause:
                first_condition = self.compare(row[table_columns.index(col1)], op1, int(val1))
                if binop:
                    second_condition = self.compare(row[table_columns.index(col2)], op2, int(val2))
                    if binop == 'AND': include = (first_condition and second_condition)
                    else: include = (first_condition or second_condition)
                else:
                    include = first_condition
            if include:
                result_row = [row[table_columns.index(col)] for col in cols]
                t.add_row(result_row)
                rows_returned += 1
        print(t)
        print('(' + str(rows_returned) + ' rows returned by query)')

if __name__=='__main__':
    cql = CQLImplementation()
    query = input('Enter a CQL SELECT query or q to quit: ')
    while query != 'q':
        cql.execute(query.upper())
        query = input('Enter a CQL SELECT query or q to quit: ')

