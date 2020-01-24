#! /usr/bin/env python
import os
from collections import deque

dir_path = os.path.dirname(os.path.realpath(__file__))

f0 = open(dir_path+"/../../../docker-volumes/node0/sgn/sgn.log", "r").readlines()
f1 = open(dir_path+"/../../../docker-volumes/node1/sgn/sgn.log", "r").readlines()
f2 = open(dir_path+"/../../../docker-volumes/node2/sgn/sgn.log", "r").readlines()

def readnext(lines, n):
    if n > len(lines) - 1:
        return "", n
    line = lines[n]
    if len(line) > 24 and line[1] != '[' and line[24] == '|':
        return line, n+1
    else:
        return readnext(lines, n+1)

def filter(log):
    queue = deque()
    n = 0
    while 1:
        line, n = readnext(log, n)
        if len(line) > 0:
            queue.append(line.rstrip('\n'))
        if n > len(log) - 1:
            return queue

def select(lines):
    time = ""
    line = ""
    index = -1
    for i in range(len(lines)):
        if len(lines[i]) < 24:
            continue
        t = lines[i][:23]
        if time == "":
            line = lines[i]
            time = t
            index = i
        else:      
            if t < time:
                line = lines[i]
                time = t
                index = i
    return line, index

def merge(logs):
    mergelog = []
    lines = []
    for log in logs:
        lines.append(log.popleft())
    while 1:
        line, index = select(lines)
        if index == -1:
            break
        mergelog.append("n%d: %s"%(index, line))
        if len(logs[index]) == 0:
            lines[index] = ""
        else:
            lines[index] = logs[index].popleft()
    return mergelog

if __name__ == '__main__':
    log0 = filter(f0)
    log1 = filter(f1)
    log2 = filter(f2)

    mergelog = merge([log0, log1, log2])
    for l in mergelog:
        print(l)