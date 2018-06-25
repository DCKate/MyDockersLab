<?php
    header('Content-Type: application/vnd.apple.mpegurl');	
    require_once __DIR__ . '/pmodules/CacheModule.php';
    session_start();
    $sessionID=session_id();
    if (isset($_SERVER['TOKEN']))
    {
        $tok = $_SERVER['TOKEN'];

        $red_ser=$_SERVER['REDIS_SERVER'];
        $red_port=$_SERVER['REDIS_PORT'];
        list($ret,$redis)=RedisConnect($red_ser, $red_port);
        if($ret==0){
            list($ret,$json)=RedisGetJson($redis,$tok);
            if($ret==0){
                $obj = json_decode($json);
                if (file_exists($obj->{'path'}))
                {
                    $content="";
                    $timestr = sprintf("%d",time());
                    $handle = fopen($obj->{'path'}, "r");
                    if ($handle) {
                        while (($line = fgets($handle)) !== false) {
                            if(preg_match("/\w.ts\s/i", $line)){
                                $ticket = hash('crc32b',$timestr.$sessionID.trim($line));
                                $content.=substr(trim($line),0,-3).".hls?ticket=".$ticket."&timestamp=".$timestr."\n";
                            }else{
                                $content.=trim($line)."\n";
                            }
                        }
                        fclose($handle);
                    }
                    print $content;
                    return;
                    
                }
            }
        }
        echo "Error !! file not found!";
        
    }
?>