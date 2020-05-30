#!/usr/bin/env python
#-*-coding: utf-8 -*-

import sys
import random
import ntpath
import os
# import numpy as np

def path_leaf(path):
    head, tail = ntpath.split(path)
    return tail or ntpath.basename(head)

def load_reversed_graph(graph):
    G = dict()
    with open(graph, 'r') as fin:
        for l in fin:
            u1, u2 = map(int, l.rstrip().split()[:2])
            if u2 not in G:
                G[u2] = list()
            G[u2].append(u1)
    return G

def constant_ic(graph, p):
    fn = os.path.splitext(path_leaf(graph))[0]
    f = open(fn+"_IC.inf","w+")
    with open(graph, 'r') as fin:
        for l in fin:
            u1, u2 = map(int, l.rstrip().split()[:2])
            f.write("%d\t%d\t%.3g\n" % (u1, u2, p))
    f.close()

def weighted_cascade(graph):
    G = load_reversed_graph(graph)
    fn = os.path.splitext(path_leaf(graph))[0]
    f = open(fn+"_WC.inf","w+")
    for u2 in G:
        for u1 in G[u2]:
            f.write("%d\t%d\t%.3g\n" % (u1, u2, 1. / len(G[u2])))
    f.close()

def tri_valency_ic(graph, p):
    fn = os.path.splitext(path_leaf(graph))[0]
    f = open(fn+"_TV.inf","w+")
    with open(graph, 'r') as fin:
        for l in fin:
            u1, u2 = map(int, l.rstrip().split()[:2])
            f.write("%d\t%d\t%.3g\n" % (u1, u2, random.choice(p)))
    f.close()

def uniform_lt(graph):
    # This is exactly the same as Weighted Cascade for the IC model
    weighted_cascade(graph)

def random_lt(graph):
    G = load_reversed_graph(graph)
    fn = os.path.splitext(path_leaf(graph))[0]
    f = open(fn+"_R.inf","w+")
    for u2 in G:
        indegree = len(G[u2])
        weights = [random.random() for a in xrange(indegree)]
        weights = [w / sum(weights) for w in weights]
        for i in xrange(indegree):
            f.write("%d\t%d\t%.3g\n" % (G[u2][i], u2, weights[i]))
    f.close()

if __name__ == '__main__':
    if len(sys.argv) < 3:
        print 'Usage: python edge_weights.py <graph> <model> [<p1> <p2> ...]'
        sys.exit(1)
    graph = sys.argv[1]
    model = int(sys.argv[2])
    p = map(float, sys.argv[3:])
    if model == 0:
        constant_ic(graph, p[0])
    elif model == 1:
        weighted_cascade(graph)
    elif model == 2:
        tri_valency_ic(graph, p)
    elif model == 3:
        uniform_lt(graph)
    elif model == 4:
        random_lt(graph)
