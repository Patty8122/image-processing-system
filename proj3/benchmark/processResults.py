import os
import csv

# Define the output directory where the files are located
mycwd = os.getcwd()
os.chdir("../editor")
output_directory = os.path.join(os.getcwd(), 'results')

# Initialize a dictionary to store 'n1' values for each file
n1_values = {}

# List of image types and parallel types
imageType = ["small", "big", "mixture"]
parallelType = ["ws", "mr"]

# Define the number of repeats
repeat = 5

# Repeats for sequential
for i in imageType:
    nameFile = f"seq_{i}.txt"
    file_path = os.path.join(output_directory, nameFile)

    n1_values[nameFile] = []

    with open(file_path, 'r') as file:
        for line in file:
            n1 = float(line.strip())
            n1_values[nameFile].append(n1)

# Repeats for parallelType
for type1 in parallelType:
    for i in imageType:
        for threadCount in [2,4,6,8,12]:
            nameFile = f"{type1}_{i}_{threadCount}.txt"
            file_path = os.path.join(output_directory, nameFile)

            n1_values[nameFile] = []

            with open(file_path, 'r') as file:
                for line in file:
                    # print(line)
                    n1 = float(line.strip())
                    n1_values[nameFile].append(n1)

# Write the gathered 'n1' values to a CSV file
csv_output_file = 'support_file.csv'

with open(csv_output_file, 'w', newline='') as csv_file:
    csv_writer = csv.writer(csv_file)

    # Write the header row
    csv_writer.writerow(['File Name', 'n1 Values'])

    # Write the data rows
    for file_name, n1_list in n1_values.items():
        csv_writer.writerow([file_name, ', '.join(map(str, n1_list))])

# print(f'CSV file "{csv_output_file}" has been created with the gathered "n1" values.')

####################################################################################################

import pandas as pd

# Load the CSV data into a Pandas DataFrame
csv_file = 'support_file.csv'
df = pd.read_csv(csv_file)

# Filter rows with 'seq' in the 'File Name' column
seq_data = df[df['File Name'].str.startswith('seq_')]

# Group the data by 'small', 'big', and 'mixture'
grouped = seq_data.groupby(seq_data['File Name'].str.extract('(small|big|mixture)')[0])

def get_average_n1_values(group):
    # Calculate the average 'n1' values for each group
    average_n1_values = group['n1 Values'].apply(lambda x: x.split(', ')[1]).astype(float).mean()

    return average_n1_values

# Calculate the average 'n1' values for each group
average_n1_values = grouped.apply(get_average_n1_values)

# convert to dictionary
average_n1_values.to_dict()
# print("The average 'n1' value for 'small' images:", average_n1_values['small'])
# print("The average 'n1' value for 'big' images:", average_n1_values['big'])
# print("The average 'n1' value for 'mixture' images:", average_n1_values['mixture'])

####################################################################################################
import pandas as pd
import matplotlib.pyplot as plt

# Load the CSV data into a Pandas DataFrame
csv_file = 'support_file.csv'
df = pd.read_csv(csv_file)

# Calculate the average 'n2' values (Time taken) for 'mr' and 'ws'
mr_data = df[df['File Name'].str.contains('mr')]
ws_data = df[df['File Name'].str.contains('ws')]
seq_data = df[df['File Name'].str.contains('seq_')]

average_time_mr = mr_data['n1 Values'].apply(lambda x: x.split(', ')[1]).astype(float).mean()
average_time_ws = ws_data['n1 Values'].apply(lambda x: x.split(', ')[1]).astype(float).mean()
average_time_seq = seq_data['n1 Values'].apply(lambda x: x.split(', ')[1]).astype(float).mean()


# print("Average Time taken for 'mr':", average_time_mr)
# print("Average Time taken for 'ws':", average_time_ws)
# print("Average Time taken for 'seq':", average_time_seq)

# Split the 'File Name' into separate columns
df[['parallelType', 'imageType', 'numThreads']] = df['File Name'].str.extract(r'([^_]+)_([^_]+)_(\d+)\.txt')

# populate the imageType column
df.loc[ df['File Name'].str.contains('seq'), 'imageType'] =  df['File Name'].apply(lambda x: x.split('.')[0][4:])
df.loc[ df['File Name'].str.contains('seq'), 'parallelType'] =  'seq'

# add an extra column for the average time taken
df['averageTime'] = df['n1 Values'].apply(lambda x: x.split(', ')[1]).astype(float)
# print(df)

# Add a new column for the speedup : seq_time / parallel_time
# based on imageType, we can calculate the speedup for each parallelType!=seq
seq_time = {
    "small": df[(df['parallelType'] == 'seq') & (df['imageType'] == 'small')]['averageTime'].values[0],
    "big": df[(df['parallelType'] == 'seq') & (df['imageType'] == 'big')]['averageTime'].values[0],
    "mixture": df[(df['parallelType'] == 'seq') & (df['imageType'] == 'mixture')]['averageTime'].values[0]
}

df['speedup'] = df.apply(lambda row: seq_time[row['imageType']] / row['averageTime'] if row['parallelType'] != 'seq' else 1, axis=1)
# print(df)
df.to_csv("speedup.csv")

# Create and plot speedup graphs
plt.figure(figsize=(10, 10))

for imageType in ['small', 'big', 'mixture']:
    plt.plot(df[(df['parallelType'] == 'ws') & (df['imageType'] == imageType)]['numThreads'],
             df[(df['parallelType'] == 'ws') & (df['imageType'] == imageType)]['speedup'],
             marker='o', label=f'ws_{imageType}')

    plt.plot(df[(df['parallelType'] == 'mr') & (df['imageType'] == imageType)]['numThreads'],
             df[(df['parallelType'] == 'mr') & (df['imageType'] == imageType)]['speedup'],
             marker='o', label=f'mr_{imageType}')
    
plt.xlabel('Number of Threads')
plt.ylabel('Speedup')
plt.title('Speedup vs. Number of Threads')
# plt.legend()
# instead of a legend, I want to add labels to each line at the end of the line
plt.legend(loc='upper left')

# Add data labels
for imageType in ['small', 'big', 'mixture']:
    ws_data = df[(df['parallelType'] == 'ws') & (df['imageType'] == imageType)]
    mr_data = df[(df['parallelType'] == 'mr') & (df['imageType'] == imageType)]

    for i in range(len(ws_data)):
        plt.text(ws_data.iloc[i]['numThreads'], ws_data.iloc[i]['speedup'], f'{ws_data.iloc[i]["speedup"]:.2f}', ha='center', va='bottom')
        plt.text(mr_data.iloc[i]['numThreads'], mr_data.iloc[i]['speedup'], f'{mr_data.iloc[i]["speedup"]:.2f}', ha='center', va='top')

# save the plot
plt.savefig('speedup.png')

# start a new figure
plt.figure(figsize=(10, 20))

import matplotlib.pyplot as plt

# Create two separate subplots for ws and mr
fig, (ax1, ax2) = plt.subplots(2, 1, figsize=(10, 10), sharex=True)

imageTypes = ['small', 'big', 'mixture']

for imageType in imageTypes:
    ws_data = df[(df['parallelType'] == 'ws') & (df['imageType'] == imageType)]
    mr_data = df[(df['parallelType'] == 'mr') & (df['imageType'] == imageType)]

    ax1.plot(ws_data['numThreads'], ws_data['speedup'], marker='o', label=f'ws_{imageType}')
    ax2.plot(mr_data['numThreads'], mr_data['speedup'], marker='o', label=f'mr_{imageType}')

ax1.set_ylabel('Speedup (ws)')
ax2.set_xlabel('Number of Threads')
ax2.set_ylabel('Speedup (mr)')

ax1.set_title('Speedup vs. Number of Threads (ws)')
ax2.set_title('Speedup vs. Number of Threads (mr)')

ax1.legend()
ax2.legend()

# Add data labels
for imageType in imageTypes:
    ws_data = df[(df['parallelType'] == 'ws') & (df['imageType'] == imageType)]
    mr_data = df[(df['parallelType'] == 'mr') & (df['imageType'] == imageType)]

    for i in range(len(ws_data)):
        ax1.text(ws_data.iloc[i]['numThreads'], ws_data.iloc[i]['speedup'], f'{ws_data.iloc[i]["speedup"]:.2f}', ha='center', va='bottom')
        ax2.text(mr_data.iloc[i]['numThreads'], mr_data.iloc[i]['speedup'], f'{mr_data.iloc[i]["speedup"]:.2f}', ha='center', va='top')

# plt.show()
plt.savefig('speedup_subplots.png')


####################################################################################################