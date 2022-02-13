<?php
namespace MapTool;
use CellTool\Cell;
//地图类
class Map {
    public $Cells=[] ;
    public $Width, $Height;
    public function loadMap($mapPath) {
        $fh = fopen($mapPath , "rb");
        fread($fh , 22);
        $this->Width = hexdec(implode("",array_reverse(str_split(bin2hex(fread($fh , 2)),2))));
        $this->Height = hexdec(implode("",array_reverse(str_split(bin2hex(fread($fh , 2)),2))));
        //回退文件头
        fseek($fh,0);
        fread($fh , 28);
        for ( $x = 0; $x < $this->Width; $x++) {
            for ($y = 0; $y < $this->Height; $y++) {
                $Cells[$x][$y] = new Cell();
            }
        }
        for ($x = 0; $x < $this->Width/2; $x++) {
            for ($y = 0; $y < $this->Height/2; $y++) {
                fread($fh , 1);
                fread($fh , 2);
            }
        }
        for ( $x = 0; $x < $this->Width; $x++) {
            for ($y = 0; $y < $this->Height; $y++) {
                $flag = bin2hex(fread($fh , 1));
                $Cells[$x][$y]->MiddleAnimationFrame = bin2hex(fread($fh , 1));
                $value = bin2hex(fread($fh , 1));
                $Cells[$x][$y]->FrontAnimationFrame = $value == 255 ? 0 : $value;
                $Cells[$x][$y]->FrontAnimationFrame &= 0x8F; 
    
                $Cells[$x][$y]->FrontFile = bin2hex(fread($fh , 1));
                $Cells[$x][$y]->MiddleFile = bin2hex(fread($fh , 1));
    
                $Cells[$x][$y]->MiddleImage = bin2hex(fread($fh , 2)) + 1;
                $Cells[$x][$y]->FrontImage = bin2hex(fread($fh , 2)) + 1;
                fread($fh , 3);
                $Cells[$x][$y]->Light = (bin2hex(fread($fh , 1)) & 0x0F) * 2;
                fread($fh , 1);
                $Cells[$x][$y]->Flag = (($flag & 0x01) != 1) || (($flag & 0x02) != 2);
            }
        }
        fclose($fh);
        
        //地图出力
        for ( $x = 0; $x < $this->Width; $x++) {
            for ($y = 0; $y < $this->Height; $y++) {
                if(!$Cells[$x][$y]->Flag) {
                    $this->Cells[$x][$y] = 1;
                    //file_put_contents("map.txt",1,FILE_APPEND);
                }else {
                    $this->Cells[$x][$y] = 0;
                    //file_put_contents("map.txt",0,FILE_APPEND);
                }
            }
            //file_put_contents("map.txt","\n",FILE_APPEND);
        }
    }
}