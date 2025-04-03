#!/bin/bash
#
#SBATCH --mail-user=divyapattisapu@cs.uchicago.edu
#SBATCH --mail-type=ALL
#SBATCH --job-name=proj1_benchmark 
#SBATCH --output=./proj1.stdout
#SBATCH --error=./proj1.stderr
#SBATCH --chdir=/home/divyapattisapu/autumn23/project-3-Patty8122/proj3/editor
#SBATCH --partition=debug 
#SBATCH --nodes=1
#SBATCH --ntasks=1
#SBATCH --cpus-per-task=12
#SBATCH --mem-per-cpu=900
#SBATCH --exclusive
#SBATCH --time=3:00:00

# module load golang/1.21.3
# module load go1.21.3 linux/amd64
# Your command here
echo $PWD
repeat=5


parallelType=("ws" "mr")
imageType=("small" "big" "mixture")

result_dir=/home/divyapattisapu/autumn23/project-3-Patty8122/proj3/editor/results

if [ ! -d "$result_dir" ]; then
    mkdir -p "$result_dir"
fi


# Repeats for sequential
for i in "${imageType[@]}"
do
    nameFile="seq_"$i".txt"
    echo "Writing to $nameFile"
    for ((j=1;j<=repeat;j++))
    do
        go run "$PWD"/editor.go $i >> "$result_dir"/$nameFile
    done
    echo "Done"
done

# Repeats for parallelType
for type1 in "${parallelType[@]}"
do 
    for i in "${imageType[@]}"
    do
        for threadCount in 2 4 6 8 12
        do
            nameFile=$type1"_"$i"_"$threadCount".txt"
            echo "Writing to $nameFile"
            for ((j=1;j<=repeat;j++))
            do
                go run "$PWD"/editor.go $i "$type1" $threadCount >> "$result_dir"/$nameFile
            done
            echo "Done"
        done
    done
done

python "$PWD"/processResults.py