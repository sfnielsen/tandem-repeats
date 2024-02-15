import pandas as pd
import matplotlib.pyplot as plt
import os

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

    # Clean and convert 'RunningTime' to numeric
    try:
        data['RunningTime'] = pd.to_numeric(data['RunningTime'].str.replace(r'[^0-9.]', ''), errors='coerce')
    except AttributeError:
        pass
    # Convert 'InputSize' to categorical for better grouping
    data['InputSize'] = pd.Categorical(data['InputSize'])
    grouped_data = data.groupby(['Algorithm', 'InputSize'])['RunningTime'].mean().reset_index()
    #Plot the data
    for algorithm in grouped_data['Algorithm'].unique():
        algorithm_data = grouped_data[grouped_data['Algorithm'] == algorithm]
        print(algorithm_data)
        plt.plot(algorithm_data['InputSize'], algorithm_data['RunningTime'], label=algorithm)
        
    plt.xlabel('Input Size')
    plt.ylabel('Running Time (ms)')
    plt.legend()
    plt.show()

def scatterplot_tandem_repeats(data:pd.DataFrame):
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
folder_path = 'time_csvs'
latest_file_path = get_latest_file(folder_path)
plot_tandem_repeats(pd.read_csv(latest_file_path, sep=","))
