import pandas as pd
import matplotlib.pyplot as plt
import os
import math
import numpy as np

#fetch the latest data in the speciale/time/csvs folder
def get_latest_file(folder_path:str):
    # Get a list of all files in the folder
    files = [os.path.join(folder_path, file) for file in os.listdir(folder_path) if os.path.isfile(os.path.join(folder_path, file))]

    # Return the latest file based on the last modification time
    return max(files, key=os.path.getmtime, default=None)

def plot_tandem_repeats(data: pd.DataFrame):
    #Group based column "Algorithm" and get the mean of the "Time" column for each "InputSize" 
        # Convert 'InputSize' to numeric
    data['InputSize'] = pd.to_numeric(data['InputSize'], errors='coerce')
    
    #convert µ to ms
    data['RunningTime'] = data['RunningTime'].apply(lambda x: float(x[:-2]) * 1000 if type(x) == str and x[-2] == 'µ' else x)
    data['RunningTime'] = data['RunningTime'].apply(lambda x: float(x[:-1]) * 1000 if type(x) == str and x[-1] == 's' else float(x))
    # Clean and convert 'RunningTime' to numeric
    try:
        data['RunningTime'] = pd.to_numeric(data['RunningTime'].str.replace(r'[^0-9.]', ''), errors='coerce')
    except AttributeError:
        pass
    # Convert 'InputSize' to categorical for better grouping
    data['InputSize'] = pd.Categorical(data['InputSize'])
    grouped_data = data.groupby(['Algorithm', 'InputSize'])['RunningTime'].mean().reset_index()
    grouped_data['StandardError'] = data.groupby(['InputSize', 'Algorithm'])['RunningTime'].sem().reset_index()['RunningTime']

    
    #Plot the data
    for algorithm in grouped_data['Algorithm'].unique():
        algorithm_data = grouped_data[grouped_data['Algorithm'] == algorithm]
        print(algorithm_data)
        plt.plot(algorithm_data['InputSize'], algorithm_data['RunningTime'], label=algorithm)
        
        # convert data to latex format for easy copy paste to latex
        #must be on the following format, with the (5.0,5.0) being the error margin for each point
        """
        coordinates {
            (0,23.1)     +- (5.0,5.0) 
            (10,27.5)    +- (10.5,10.5)
            (20,32)      +- (7.5,7.5)
            (30,37.8)    +- (12.0,12.0)
            (40,44.6)    +- (7.0,7.0)
            (60,61.8)    +- (8.0,8.0)
            (80,83.8)    +- (16.0,16.0)
            (100,114)    +- (14.0,14.0)
        };
        """
        print(f"coordinates {{")
        for i in range(len(algorithm_data['InputSize'])):
            # upper lower bound is the +- part
            print(f"({algorithm_data['InputSize'].iloc[i]},{round(algorithm_data['RunningTime'].iloc[i],2)}) +- ({round(2*algorithm_data['StandardError'].iloc[i],2)},{round(2*algorithm_data['StandardError'].iloc[i],2)})")

        print(f"}};")

        
    plt.xlabel('Input Size')
    plt.ylabel('Running Time (ms)')
    plt.legend()
    plt.show()
    
    

def scatterplot_tandem_repeats(data:pd.DataFrame):
    data['InputSize'] = pd.to_numeric(data['InputSize'], errors='coerce')

    # convert runnningtime to numeric, for example 's' means seconds so convert it to ms
    data['RunningTime'] = data['RunningTime'].apply(lambda x: float(x[:-2]) * 1000 if type(x) == str and x[-2] == 'µ' else x)
    data['RunningTime'] = data['RunningTime'].apply(lambda x: float(x[:-1]) * 1000 if type(x) == str and x[-1] == 's' else float(x))
    # Convert 'InputSize' to numeric
    data['InputSize'] = pd.to_numeric(data['InputSize'], errors='coerce')

    # Clean and convert 'RunningTime' to numeric
    try:
        data['RunningTime'] = pd.to_numeric(data['RunningTime'].str.replace(r'[^0-9.]', ''), errors='coerce')
    except AttributeError:
        pass

    # Plot all data points
    for algorithm in data['Algorithm'].unique():
        algorithm_data = data[data['Algorithm'] == algorithm]
        plt.scatter(algorithm_data['InputSize'], algorithm_data['RunningTime'], label=algorithm, marker='x')

    plt.xlabel('Input Size')
    plt.ylabel('Running Time (ms)')
    plt.legend()
    plt.show()

def plot_expected_time_complexity_test(data:pd.DataFrame):

    
    data['InputSize'] = pd.to_numeric(data['InputSize'], errors='coerce')

    # convert runnningtime to numeric, for example 's' means seconds so convert it to ms
    data['RunningTime'] = data['RunningTime'].apply(lambda x: float(x[:-2]) * 1000 if type(x) == str and x[-2] == 'µ' else x)
    data['RunningTime'] = data['RunningTime'].apply(lambda x: float(x[:-1]) * 1000 if type(x) == str and x[-1] == 's' else float(x))

    # Clean and convert 'RunningTime' to numeric
    try:
        data['RunningTime'] = pd.to_numeric(data['RunningTime'].str.replace(r'[^0-9.]', ''), errors='coerce')
    except AttributeError:
        pass

    # Group by InputSize and Algorithm and calculate the average RunningTime for each group
    grouped_df = data.groupby(['InputSize', 'Algorithm']).agg({'RunningTime': 'mean','Complexity': 'first'}).reset_index()
    # Plot the data with RunningTime divided by InputSize
    fig, ax = plt.subplots()


    # Calculate the standard error for each group and add a column for lower and upper part of the mean +- 2*standard error
    grouped_df['StandardError'] = data.groupby(['InputSize', 'Algorithm'])['RunningTime'].sem().reset_index()['RunningTime']
    grouped_df['LowerBound'] = grouped_df['RunningTime'] - 2 * grouped_df['StandardError']
    grouped_df['UpperBound'] = grouped_df['RunningTime'] + 2 * grouped_df['StandardError']
    print(grouped_df)
    # Iterate over unique algorithms and plot them
    for algorithm in grouped_df['Algorithm'].unique():


        algorithm_data = grouped_df[grouped_df['Algorithm'] == algorithm]
        algorithm_data['InputSize']
        selected_complexity = algorithm_data['Complexity'].iloc[0]

        if selected_complexity == 'nlogn':
            complexity_function = lambda x: x / ((algorithm_data['InputSize'])*np.log2( algorithm_data['InputSize']))
        elif selected_complexity == 'n':
            complexity_function = lambda x: x / algorithm_data['InputSize']
        elif selected_complexity == 'n^2':
            #algorithm_data['InputSize'] /= 1000
            complexity_function = lambda x: x / (algorithm_data['InputSize'] ** 2)
        else:
            complexity_function = lambda x: x
        print("YASONDOASND")
        algorithm_data['RunningTime'] = complexity_function(algorithm_data['RunningTime'])
        print(f"coordinates {{")
        for i in range(len(algorithm_data['InputSize'])):
            # upper lower bound is the +- part
            print(f"({algorithm_data['InputSize'].iloc[i]},{round(algorithm_data['RunningTime'].iloc[i],6)})")

        print(f"}};")
        ax.plot(algorithm_data['InputSize'], algorithm_data['RunningTime'], label=f'{algorithm} ({selected_complexity})')

        ax.set_xlabel('InputSize')
        ax.set_ylabel('log(RunningTime / InputSize)')
    ax.legend()
    plt.show()

## create a similar plot but with scatter plot, so no values are aggregated
def plot_expected_time_complexity_test_scatter(data:pd.DataFrame):
    data['InputSize'] = pd.to_numeric(data['InputSize'], errors='coerce')

    # convert runnningtime to numeric, for example 's' means seconds so convert it to ms
    data['RunningTime'] = data['RunningTime'].apply(lambda x: float(x[:-2]) * 1000 if type(x) == str and x[-2] == 'µ' else x)
    data['RunningTime'] = data['RunningTime'].apply(lambda x: float(x[:-1]) * 1000 if type(x) == str and x[-1] == 's' else float(x))

    # Clean and convert 'RunningTime' to numeric
    try:
        data['RunningTime'] = pd.to_numeric(data['RunningTime'].str.replace(r'[^0-9.]', ''), errors='coerce')
    except AttributeError:
        pass

    # Group by InputSize and Algorithm and calculate the average RunningTime for each group
    grouped_df = data
    # Plot the data with RunningTime divided by InputSize
    fig, ax = plt.subplots()


    # Calculate the standard error for each group and add a column for lower and upper part of the mean +- 2*standard error
    grouped_df['StandardError'] = data.groupby(['InputSize', 'Algorithm'])['RunningTime'].sem().reset_index()['RunningTime']
    grouped_df['LowerBound'] = grouped_df['RunningTime'] - 2 * grouped_df['StandardError']
    grouped_df['UpperBound'] = grouped_df['RunningTime'] + 2 * grouped_df['StandardError']
    print(grouped_df)
    # Iterate over unique algorithms and plot them
    for algorithm in grouped_df['Algorithm'].unique():


        algorithm_data = grouped_df[grouped_df['Algorithm'] == algorithm]
        algorithm_data['InputSize']
        selected_complexity = algorithm_data['Complexity'].iloc[0]

        if selected_complexity == 'nlogn':
            complexity_function = lambda x: x / algorithm_data['InputSize']
            #complexity_function = lambda x: x / ((algorithm_data['InputSize'])*np.log2( algorithm_data['InputSize']))
        elif selected_complexity == 'n':
            complexity_function = lambda x: x / algorithm_data['InputSize']
        elif selected_complexity == 'n^2':
            algorithm_data
            complexity_function = lambda x: x / (algorithm_data['InputSize'] ** 2)
        else:
            complexity_function = lambda x: x
        print("YASONDOASND")
        algorithm_data['RunningTime'] = complexity_function(algorithm_data['RunningTime'])

        print(f"coordinates {{")
        for i in range(len(algorithm_data['InputSize'])):
            # upper lower bound is the +- part
            print(f"({algorithm_data['InputSize'].iloc[i]},{round(algorithm_data['RunningTime'].iloc[i],6)})")

        print(f"}};")
        ax.scatter(algorithm_data['InputSize'], algorithm_data['RunningTime'], label=f'{algorithm} ({selected_complexity})')

        ax.set_xlabel('InputSize')
        ax.set_ylabel('log(RunningTime / InputSize)')
    ax.legend()
    plt.show()


folder_path = 'time_csvs'
latest_file_path = get_latest_file(folder_path)
print(latest_file_path)
plot_expected_time_complexity_test_scatter(pd.read_csv(latest_file_path, sep=","))
