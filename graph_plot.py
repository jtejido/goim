import matplotlib.pyplot as plt
  
x = []
y = []
x1 = []
y1 = []

x2 = []
y2 = []

file1 = open('output/output_graph/twitter_combined_WC_discountdegree_15_50_1487723611282_1682961523703.txt', 'r')
Lines = file1.readlines()

count = 0
for line in Lines:
    count += 1
    line = line.strip()
    line = line[0:len(line)-2]
    lines = [i for i in line.split(", ")]
    if len(lines) == 1:
        continue

    x2.append(count)
    y2.append(lines[0])

    for i in range(len(lines)):
        x.append(count)
        y.append(int(lines[i]))
    x1.append(x)
    y1.append(y)

plt.scatter(x1, y1)
plt.show()

plt.plot(x2, y2)
plt.show()
