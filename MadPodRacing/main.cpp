#include <math.h>

#include <algorithm>
#include <iostream>
#include <string>
#include <vector>

using namespace std;

/**
 * Auto-generated code below aims at helping you parse
 * the standard input according to the problem statement.
 **/
struct Checkpoint {
  int x, y;
  bool operator==(const Checkpoint &a) {
    if (x == a.x && y == a.y) return true;
    return false;
  }
  bool operator!=(const Checkpoint &a) { return !(*this == a); }
  double dist(const Checkpoint &a) const {
    return pow(a.x - x, 2) + pow(a.y - y, 2);
  }
};

void output_boost(pair<int, int> coordinate) {
  cout << coordinate.first << " " << coordinate.second << " BOOST" << endl;
}
void output(pair<int, int> coordinate) {
  cout << coordinate.first << " " << coordinate.second << " 100" << endl;
}

int boost_timing = -1;
int calc_boost_timing(const vector<Checkpoint> &checkpoints) {
  if (boost_timing != -1) return boost_timing;
  vector<double> dists;
  for (int i = 0; i < checkpoints.size(); i++) {
    int compare_idx = i - 1;
    if (i == 0) compare_idx = checkpoints.size() - 1;
    dists.push_back(checkpoints[i].dist(checkpoints[compare_idx]));
  }
  return max_element(dists.begin(), dists.end()) - dists.begin();
}

int regest_checkpoint(vector<Checkpoint> &checkpoints,
                      const Checkpoint &next_checkpoint) {
  for (int i = 0; i < checkpoints.size(); i++) {
    if (checkpoints[i] == next_checkpoint) {
      return i;
    }
  }
  checkpoints.push_back(next_checkpoint);
  return checkpoints.size() - 1;
}

int lap = 0;
int pre_target_no = 0;
int calc_lap_no(int target_no) {
  if (pre_target_no == target_no) return lap;
  pre_target_no = target_no;

  if (target_no == 0) lap++;
  return lap;
}

int gcd(int a, int b) {
  if (a % b == 0) return b;
  return gcd(b, a % b);
}

// Fx : Fy の比を返す
pair<int, int> calc_power_balance(double dist_x, double dist_y, double v_y0,
                                  double v_x0) {
  int child = v_x0 * dist_y + v_y0 * dist_x;
  int mother = v_x0 * v_y0;
  if (child == 0 && mother == 0) return {0, 0};
  if (child == 0) return {100, 0};
  if (mother == 0) return {0, 100};
  int g = gcd(child, mother);
  return {mother / g, child / g};
}

int pre_x = 1e9, pre_y = 1e9;
pair<int, int> calc_target_coordinate(int x, int y, int target_x,
                                      int target_y) {
  int dist_x = target_x - x;
  int dist_y = target_y - y;
  int speed_x = pre_x == 1e9 ? 0 : pre_x - x;
  int speed_y = pre_y == 1e9 ? 0 : pre_y - y;
  pre_x = x;
  pre_y = y;
  pair<int, int> p = calc_power_balance(dist_x, dist_y, speed_x, speed_y);
  cerr << p.first << " " << p.second << " " << x << " " << y << endl;
  return {x + p.first, y + p.second};
}

int main() {
  vector<Checkpoint> checkpoints;
  // game loop
  while (1) {
    int x, y;
    Checkpoint next_checkpoint;
    int next_checkpoint_dist;
    // angle between your pod orientation and the direction of the next
    int next_checkpoint_angle;

    cin >> x >> y >> next_checkpoint.x >> next_checkpoint.y >>
        next_checkpoint_dist >> next_checkpoint_angle;
    cin.ignore();
    int opponent_x, opponent_y;
    cin >> opponent_x >> opponent_y;
    cin.ignore();

    int target_no = regest_checkpoint(checkpoints, next_checkpoint);
    int lap_no = calc_lap_no(target_no);

    int boost_timing = -1;
    if (lap_no > 0) {
      boost_timing = calc_boost_timing(checkpoints);
    }

    // Write an action using cout. DON'T FORGET THE "<< endl"
    // To debug: cerr << "Debug messages..." << endl;

    // You have to output the target position
    // followed by the power (0 <= thrust <= 100)
    // i.e.: "x y thrust"
    pair<int, int> target_coordinate =
        calc_target_coordinate(x, y, next_checkpoint.x, next_checkpoint.y);
    if (lap_no == 1 && boost_timing == target_no &&
        abs(next_checkpoint_angle) < 30) {
      output_boost(target_coordinate);
    } else {
      output(target_coordinate);
    }
  }
}