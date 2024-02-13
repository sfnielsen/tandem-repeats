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
    #Group based column "Algorithm" and get the mean of the "Time" column for each "input" value
    grouped_data = data.groupby(['Algorithm', 'InputSize']).mean().reset_index()
    #Plot the data
    for algorithm in grouped_data['Algorithm'].unique():
        algorithm_data = grouped_data[grouped_data['Algorithm'] == algorithm]
        plt.plot(algorithm_data['InputSize'], algorithm_data['RunningTime'], label=algorithm)
    plt.xlabel('Input Size')
    plt.ylabel('Running Time (ms)')
    plt.show()

folder_path = 'speciale/time_csvs'
latest_file_path = get_latest_file(folder_path)
plot_tandem_repeats(pd.read_csv(latest_file_path))
