#include <iostream>
#include <string>
#include <vector>
#include <algorithm>
#include <vector>
#include <math.h>

using namespace std;

/**
 * Auto-generated code below aims at helping you parse
 * the standard input according to the problem statement.
 **/
struct Checkpoint{
    int x, y;
    bool operator==(const Checkpoint &a){
        if(x == a.x && y == a.y)return true;
        return false;
    }
    bool operator!=(const Checkpoint &a){
        return !(*this == a);
    }
    double dist(const Checkpoint &a) const{
        return pow(a.x - x, 2) + pow(a.y - y, 2);
    }
};

void output_boost(Checkpoint checkpoint){
    cout << checkpoint.x << " " << checkpoint.y << " BOOST" << endl;
}
void output(Checkpoint checkpoint, int thrust){
    cout << checkpoint.x << " " << checkpoint.y << " " << thrust << endl;
}



int boost_timing = -1;
int calc_boost_timing(const vector<Checkpoint> &checkpoints){
    if(boost_timing != -1)return boost_timing;
    vector<double> dists;
    for(int i=0;i<checkpoints.size();i++){
        int compare_idx = i-1;
        if(i == 0)compare_idx = checkpoints.size() - 1;
        dists.push_back(checkpoints[i].dist(checkpoints[compare_idx]));
    }
    return max_element(dists.begin(), dists.end()) - dists.begin();
}

int regest_checkpoint(vector<Checkpoint> &checkpoints, const Checkpoint &next_checkpoint){
    for(int i=0;i<checkpoints.size();i++){
        if(checkpoints[i] == next_checkpoint){
            return i;
        }
    }
    checkpoints.push_back(next_checkpoint);
    return checkpoints.size() - 1;
}

int lap = 0;
int pre_target_no = 0;
int calc_lap_no(int target_no){
    if(pre_target_no == target_no)return lap;
    pre_target_no = target_no;

    if(target_no == 0)lap++;
    return lap;
}

int calc_thrust2(int next_checkpoint_angle){
    // 0 ~ 1 
    double angle_closeness = 1 - abs((double)next_checkpoint_angle) / 180;
    
    if(angle_closeness > 0.5)return 100;
    return 0;
}


int calc_thrust(int next_checkpoint_angle){
    // 0 ~ 1 
    double angle_closeness = 1 - abs((double)next_checkpoint_angle) / 180;
    
    auto sigmoid = [](double x){
        double a = 10;
        return 1 / (1 + pow(M_E, -a * x));
    };

    double thrust_n = sigmoid(angle_closeness - 0.5);

    int res = 100 * thrust_n;
    if(res == 99)return 100;
    return res;
}


int main()
{
    vector<Checkpoint> checkpoints;
    // game loop
    while (1) {
        int x;
        int y;
        Checkpoint next_checkpoint;
        int next_checkpoint_dist; // distance to the next checkpoint
        int next_checkpoint_angle; // angle between your pod orientation and the direction of the next checkpoint
        cin >> x >> y >> next_checkpoint.x >> next_checkpoint.y >> next_checkpoint_dist >> next_checkpoint_angle; cin.ignore();
        int opponent_x;
        int opponent_y;
        cin >> opponent_x >> opponent_y; cin.ignore();

        int target_no = regest_checkpoint(checkpoints, next_checkpoint);
        int lap_no = calc_lap_no(target_no);

        int boost_timing = -1;
        if(lap_no > 0){
            boost_timing = calc_boost_timing(checkpoints);
        }

        // Write an action using cout. DON'T FORGET THE "<< endl"
        // To debug: cerr << "Debug messages..." << endl;


        // You have to output the target position
        // followed by the power (0 <= thrust <= 100)
        // i.e.: "x y thrust"
        int thrust = calc_thrust(next_checkpoint_angle);
        cerr << "lap_no" << lap_no << " target_no" << target_no << " boost_timing" << boost_timing << endl;
        if(lap_no == 1 && boost_timing == target_no && abs(next_checkpoint_angle) < 30){
            output_boost(next_checkpoint);
        }else{
            output(next_checkpoint, thrust);
        }
    }
}