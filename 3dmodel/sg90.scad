//精度
$fs=0.5;
$fa=0.5;

// 长22.6mm  宽12.4 高15.5，线宽3.5mm
//孔 5mm
space=0.5;
sgLength=22.6 + space;
sgWidth=12.6 + space;
sgHeight=15.5 + space;
lineWidth=3.5;
board=1;
cube([sgLength + board *2, sgWidth + board * 2, 1]);
cube([sgLength + board *2, 1, sgHeight + board * 2]);
translate([0, sgWidth + board, 0]) {
    cube([sgLength + board *2, 1, sgHeight + board * 2]);
}
cube([1, sgWidth + board * 2, sgHeight + board * 2]);

difference(){
    translate([sgLength + board, 0, 0]) {
        cube([1, sgWidth + board * 2, sgHeight + board * 2]);
    }
    translate([sgLength, (sgWidth - lineWidth)/2 + board, 0]) {
        cube([3, lineWidth, sgHeight + board * 3]);
    }
}



bracketSpace=2;//支撑间隔
bracketHeight=15;//支撑高度
bracketWidth=sgWidth + board * 2;//支撑宽度
bracketHole=5/2;//孔半径
holeHeight=6;//孔到盒子距离
transX=(sgLength - bracketSpace)/2;
translate([transX + 0, 0, -bracketHeight]) {
    difference() {
        translate([0, sgWidth/2 + board, 12]) {
            cube([board, bracketWidth, bracketHeight/2], center=true);
        }
        translate([0, sgWidth/2 + board, holeHeight]) rotate([0, 90, 0]) {
            cylinder(h=10, r1= bracketHole, r2= bracketHole, center=true, cy= 1);
        }
    }
    difference() {
        translate([0, sgWidth/2 + board, bracketHeight/2]) rotate([0, 90, 0]) {
            cylinder(h=board, r1= bracketWidth/2, r2= bracketWidth/2, center=true, cy= 1);
        }
        translate([0, sgWidth/2 + board, holeHeight]) rotate([0, 90, 0]) {
            cylinder(h=10, r1= bracketHole, r2= bracketHole, center=true, cy= 1);
        }
    }
}
translate([transX + bracketSpace + board, 0, -bracketHeight]) {
     difference() {
        translate([0, sgWidth/2 + board, 12]) {
            cube([board, bracketWidth, bracketHeight/2], center=true);
        }
        translate([0, sgWidth/2 + board, holeHeight]) rotate([0, 90, 0]) {
            cylinder(h=10, r1= bracketHole, r2= bracketHole, center=true, cy= 1);
        }
    }
    difference() {
        translate([0, sgWidth/2 + board, bracketHeight/2]) rotate([0, 90, 0]) {
            cylinder(h=board, r1= bracketWidth/2, r2= bracketWidth/2, center=true, cy= 1);
        }
        translate([0, sgWidth/2 + board, holeHeight]) rotate([0, 90, 0]) {
            cylinder(h=10, r1= bracketHole, r2= bracketHole, center=true, cy= 1);
        }
    }
}
