filename = 'C:/Users/Uhini Mukherjee/Desktop/goim/output/twitter_combined_WC_discountdegree_15_50_1487723611282_1682961523703.log'

fp = open(filename,'r')
data = []
write = open("C:/Users/Uhini Mukherjee/Desktop/goim/output/output_graph/twitter_combined_WC_discountdegree_15_50_1487723611282_1682961523703.txt", 'a')
write.write("\n")
line = fp.readline()
while len(line) > 0:
    line = line.split("[")[1][:-2]
    # print(line)
    write.write(line+"\n")
    line = fp.readline()

write.close()

