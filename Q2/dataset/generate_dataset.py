import random
import argparse

def generate_knn_dataset(num_points):
    output_file = "./dataset/dataset.txt"
    with open(output_file, 'w') as f:
        for _ in range(num_points):
            x = round(random.uniform(-100, 100), 2)
            y = round(random.uniform(-100, 100), 2)
            f.write(f"{x} {y}\n")
    
    print(f"Dataset generated and saved to {output_file}")

if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument('--points', type=int, required=True, help="Number of data points to generate.")

    args = parser.parse_args()

    generate_knn_dataset(args.points)
