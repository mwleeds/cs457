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
                columns = [tablename + str(i+1) for i in range(len(lines[0].split()))] + [tablename + '.writetime']
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
        # ideally this method would be broken into smaller methods but it seems to work
        # this regex matches a SELECT query on tables A, B, or C with optional JOIN and WHERE clauses
        select_re = r'^SELECT\s+(?P<cols>.+?)\s+' + \
                    r'FROM\s+(?P<table>[A-C])' + \
                    r'((\s+JOIN\s+(?P<join_table1>[A-C])\s+ON\s+(?P<join_condition1>.+?))?' + \
                    r'(\s+JOIN\s+(?P<join_table2>[A-C])\s+ON\s+(?P<join_condition2>.+?))?' + \
                    r'(\s+WHERE\s+(?P<where_clause>.+?))?)?;$'
        tokens = re.match(select_re, query)
        if tokens == None:
            raise Exception('Invalid query entered')
        # get tokens
        cols = tokens.group('cols')
        table = tokens.group('table')
        join_table1 = tokens.group('join_table1')
        join_condition1 = tokens.group('join_condition1')
        join_table2 = tokens.group('join_table2')
        join_condition2 = tokens.group('join_condition2')
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
        if join_condition1:
            join_condition_re = r'(?P<left_col>.+?)=(?P<right_col>.+?)$'
            join_condition_tokens = re.match(join_condition_re, join_condition1)
            join1_left_col = join_condition_tokens.group('left_col')
            join1_right_col = join_condition_tokens.group('right_col')
        if join_condition2:
            join_condition_re = r'(?P<left_col>.+?)=(?P<right_col>.+?)$'
            join_condition_tokens = re.match(join_condition_re, join_condition2)
            join2_left_col = join_condition_tokens.group('left_col')
            join2_right_col = join_condition_tokens.group('right_col')
        # expand the * for columns if necessary
        table_columns = self.tables[table]['columns'][:] # make a copy of the columns list
        if join_table1: table_columns += self.tables[join_table1]['columns']
        if join_table2: table_columns += self.tables[join_table2]['columns']
        if cols == '*': cols = table_columns
        else: cols = cols.replace(' ', '').split(',')
        t = PrettyTable(cols)
        rows_returned = 0
        raw_rows = []
        # execute the first JOIN if necessary
        if join_table1:
            left_table, right_table = join1_left_col[0], join1_right_col[0]
            # for better performance this could be a hash join. here it's O(n^2)
            for left_row in self.tables[left_table]['rows']:
                for right_row in self.tables[right_table]['rows']:
                    if left_row[self.tables[left_table]['columns'].index(join1_left_col)] == right_row[self.tables[right_table]['columns'].index(join1_right_col)]:
                        raw_rows.append(left_row + right_row)
        else: # no JOIN operation was specified
            raw_rows = self.tables[table]['rows']
        # execute the second JOIN if necessary
        if join_table2:
            # if the JOIN condition's columns were specified in the wrong order, fix it
            right_col = join2_right_col if join_table2 in join2_right_col else join2_left_col
            left_col = join2_left_col if join_table2 not in join2_left_col else join2_right_col
            new_raw_rows = []
            for left_row in raw_rows:
                for right_row in self.tables[join_table2]['rows']:
                    if left_row[table_columns.index(left_col)] == right_row[self.tables[join_table2]['columns'].index(right_col)]:
                        new_raw_rows.append(left_row + right_row)
            raw_rows = new_raw_rows
        # enforce the where clause and pare down the columns
        for row in raw_rows:
            include = True
            # for better performance this could be executed before the JOIN
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

