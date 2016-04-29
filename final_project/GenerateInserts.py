#!/usr/bin/python3

writetime = 0

with open('A.txt') as f:
    for line in f:
        writetime += 1
        print('INSERT INTO A (A1, A2, writetime) VALUES (' + ', '.join(line.split()) + ', ' + str(writetime) + ');')

print()

with open('B.txt') as f:
    for line in f:
        writetime += 1
        print('INSERT INTO B (B1, B2, B3, writetime) VALUES (' + ', '.join(line.split()) + ', ' + str(writetime) + ');')

print()

with open('C.txt') as f:
    for line in f:
        writetime += 1
        print('INSERT INTO C (C1, C2, C3, C4, writetime) VALUES (' + ', '.join(line.split()) + ', ' + str(writetime) + ');')

