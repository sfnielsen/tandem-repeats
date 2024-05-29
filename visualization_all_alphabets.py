import pandas as pd
import matplotlib.pyplot as plt
import os
import math
import numpy as np


alphabet_label_mapper = {
    'AB': 'ab',
    'A': 'a',
    'ACGT':'DNA',
    'ACDEFGHIKLMNPQRSTVWY':'Protein',
    'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789':'ASCII',
    'Byte':'Byte',
    'fib':'Fibonacci',
    'max':'Max',
}
def get_latest_file(folder_path:str):
    # Get a list of all files in the folder
    files = [os.path.join(folder_path, file) for file in os.listdir(folder_path) if os.path.isfile(os.path.join(folder_path, file))]

    # Return the latest file based on the last modification time
    return max(files, key=os.path.getmtime, default=None)

def plot_tandem_repeats_all_alphabet(data: pd.DataFrame):
    #Group based column "Algorithm" and get the mean of the "Time" column for each "InputSize" 
    # Convert 'InputSize' to numeric
    data['InputSize'] = pd.to_numeric(data['InputSize'], errors='coerce')
    data['RunningTime'] = data['RunningTime'].apply(lambda x: float(x[:-2]) * 1000 if type(x) == str and x[-2] == 'µ' else x)
    data['RunningTime'] = data['RunningTime'].apply(lambda x: float(x[:-1]) * 1000 if type(x) == str and x[-1] == 's' else float(x))
    #data['RunningTime'] =  np.log(data['RunningTime'])
    # Clean and convert 'RunningTime' to numeric
    try:
        data['RunningTime'] = pd.to_numeric(data['RunningTime'].str.replace(r'[^0-9.]', ''), errors='coerce')
    except AttributeError:
        pass
    # Convert 'InputSize' to categorical for better grouping
    data['InputSize'] = pd.Categorical(data['InputSize'])
    grouped_data = data.groupby(['Algorithm','Alphabet', 'InputSize'])['RunningTime'].mean().reset_index()
    grouped_data['StandardError'] = data.groupby(['InputSize', 'Alphabet', 'Algorithm'])['RunningTime'].sem().reset_index()['RunningTime']
    print(grouped_data)
    
    #Plot the data
    for algorithm in grouped_data['Algorithm'].unique():
        for alphabet in grouped_data['Alphabet'].unique():
            algorithm_data = grouped_data[(grouped_data['Alphabet'] == alphabet) & (grouped_data['Algorithm'] == algorithm)]

            # Remove first 4 rows of series
            algorithm_data = algorithm_data.iloc[2:]
            plt.plot(algorithm_data['InputSize'],algorithm_data['RunningTime'], label=alphabet_label_mapper[alphabet])

            # convert data to latex format for easy copy paste to latex
            #must be on the following format, with the (5.0,5.0) being the error margin for each point
            print(alphabet)
            print(f"coordinates {{")
            for i in range(len(algorithm_data['InputSize'])):
                # upper lower bound is the +- part
                print(f"({algorithm_data['InputSize'].iloc[i]},{round(algorithm_data['RunningTime'].iloc[i],2)}) +- ({round(2*algorithm_data['StandardError'].iloc[i],2)},{round(2*algorithm_data['StandardError'].iloc[i],2)})")

            print(f"}};")
            print("")

        
    plt.xlabel('Input Size')
    plt.ylabel('Running Time (ms)')
    plt.legend()
    plt.show()



folder_path = 'time_csvs'
latest_file_path = get_latest_file(folder_path)
print(latest_file_path)
plot_tandem_repeats_all_alphabet(pd.read_csv(latest_file_path, sep=","))