from folder_to_norm_latencies import extract_norm_latencies
from extract_latencies import extract_latencies
from latencies_to_csv import latencies_to_csv
from csvs_to_plot import cdf_csvs_to_plot, tput_wp_plot
import os
from pathlib import Path
import sys
import subprocess
import json
import numpy as np



########## PRIMARY PLOTTING CODE ############
# This Code plots figures 4,5,6, and 9 and gives directions how to produce other figures


# File path has structure: TIMESTAMP / FIG# / PROTOCOL/ CLIENT
# Make sure to run when current wording directory is plot_figs/
def main(results_path):

    plot_target_directory = Path("plots")
    csv_target_directory = Path("csvs")

    print("Plotting...")
    plot_fig5(plot_target_directory, csv_target_directory, results_path)

def plot_fig5(plot_target_directory, csv_target_directory, latency_folder):
    read_csvs, write_csvs, _, _ = calculate_csvs_cdf("5", csv_target_directory, latency_folder)

    # Reads
    cdf_csvs_to_plot(plot_target_directory, "5", read_csvs, is_for_reads=True)

    # Writes
    cdf_csvs_to_plot(plot_target_directory, "5" + "-write", write_csvs, is_for_reads=False)    

# Returns a tuple of tuple of csv paths.
# This is used for figs 4 , 5 and 9
def calculate_csvs_cdf(figure_name, csv_target_directory,latency_folder):
    write_latencies = extract_norm_latencies(latency_folder, is_for_reads=False)
    read_latencies = extract_norm_latencies(latency_folder, is_for_reads=True)
    read_csvs = {}
    write_csvs = {}

    read_log_csvs = {}
    write_log_csvs = {}

    # read
    norm_cdf_csv, norm_log_cdf_csv= latencies_to_csv(csv_target_directory, read_latencies, figure_name)
    read_csvs["pineapple"] = norm_cdf_csv
    read_log_csvs["pineapple"] = norm_log_cdf_csv

    # write
    norm_cdf_csv, norm_log_cdf_csv = latencies_to_csv(csv_target_directory, write_latencies, figure_name + "-write")
    write_csvs["pineapple"] = norm_cdf_csv
    write_log_csvs["pineapple"] = norm_log_cdf_csv

    return read_csvs, write_csvs, read_log_csvs, write_log_csvs

    
# Delete and fix packaging
def check_cmd_output(cmd):
   # output = subprocess.check_output(cmd)
    ps = subprocess.Popen(cmd,shell=True,stdout=subprocess.PIPE,stderr=subprocess.STDOUT)
    output = ps.communicate()[0]
    return output.decode("utf-8").strip("\n") 


# returns newest results. Assumes results are in ../results
def most_recent_results():
    results_dir = "../results/"
    return results_dir + check_cmd_output("ls " +  results_dir + "| sort -r | head -n 1")


def usage():
    print("Usage: python3 plot_figs.py RESULTS_PATH")

if __name__ == "__main__":
    l = len(sys.argv)
    if l == 1:
        main(most_recent_results())
    elif l == 2:
        main(sys.argv[1])
    else :
        usage()

    
