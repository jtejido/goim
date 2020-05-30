# goim
Influence Maximization in Go ([CELF][1], [TIM][2], MaxDiscount, DegreeDiscount, Pruned Monte-Carlo)

## Objective

We have a project which requires us to find, among a list of writers/bloggers/sponsored celebrities 
and their co-authorship/co-working with one another, a shortlist of people that can effectively market our 
product (say, a phone) to the most number of people.

At first, it may seem as simple as counting all the blogs/likes/re-tweets acquired by unique users,
and finding who among them obtain the most number of media engagements. 

But this problem should not be tackled at face value.

It appears that the more connections people have, the more effective they are at marketing a product. 
A person may have a small number of posts/tweets/published work online, but given their vast connections, 
they can be more effective than the most hard-working personality.

## Influence Maximization (IM)

Word-of-mouth or viral marketing is believed to be a promising marketing strategy. 
The increasing popularity of online social networks provides a good tool to enable large scale viral marketing.

The problem of exhausting or extracting people who can maximize revenue by promoting newly-launched 
products to his/her connections (peers, relative, friends, followers, etc.) is the main goal.

**Influence Maximization** is *[NP-hard][3]*.

It turns out that the deterministic influence maximization problem is NP-hard.

To prove it is NP-hard, one must find an example of an NP-hard problem that can be reduced to influence maximization.

This would show that if we had a solution to the influence maximization problem, we could easily translate that into a solution to a problem that is known to be hard, which means influence maximization is a hard problem (e.g., IC reduced from k-max cover problem and LT reduced from vertex cover problem).

Given a promotion budget, the goal of IM is to select the best k-seed nodes from an influence graph.
An influence graph is essentially a graph with influence probabilities among nodes representing social network users.

Given a social network, a diffusion model with given parameters, and a number 𝑘,
find a seed set 𝑆 (or people) of at most 𝑘-nodes (number of people) such that the influence spread of 𝑆 is maximized.

## Possible Use Cases

1. Who are the most effective at influencing others with fake news.
2. Effective Election campaigns.
3. Viral Spread.
4. Others described [here][4].


## Parameters

See the **config.toml** file for parameter options needed to run the evaluator or run ./goim -h for help.

```bash
$ ./goim -h
  -algorithm string
        Seed-selection algorithm. (default "pmc")
  -conf string
        config file location (default "config.toml")
  -cpuprofile string
        write cpu profile to location
  -graph string
        Path of graph file. (default "graphs/hep_IC_0.1.inf")
  -log string
        write log to location
  -model string
        Diffusion model to use. (default "ic")
  -output string
        Path for output files. (default "output")
  -seed int
        Seed of rng. (default 1487723611282)
  -seeds int
        Number of seeds in each trial. (default 25)
  -trials int
        Number of trials. (default 1)
```


[1]: <http://snap.stanford.edu/class/cs224w-readings/goyal11celf.pdf> "A. Goyal, W. Lu, L. Lakshmanan. CELF++: Optimizing the Greedy Algorithm for Influence Maximization in Social Networks. WWW 2011"

[2]: <http://arxiv.org/pdf/1404.0900v2.pdf> "Y. Tang, X. Xiao, and Y. Shi. Influence maximization: Near-optimal time complexity meets practical efficiency. SIGMOD 2014"

[3]: <https://www.cs.cornell.edu/home/kleinber/kdd03-inf.pdf> "D. Kempe, J. Kleinberg, E. Tardos. Maximizing the Spread of Influence through a Social Network."

[4]: <https://dl.acm.org/doi/pdf/10.1145/2503792.2503797> "A. Guille, H. Hacid, C. Favre, and D. A. Zighed, Information diffusion in online social networks: A survey. SIGMOD 2013."
