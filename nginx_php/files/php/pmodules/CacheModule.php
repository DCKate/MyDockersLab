<?php
    require_once __DIR__ . '/CheckModule.php';
    function RedisConnect($addr,$port){
        $redis = new Redis();
        $redis->connect($addr,$port);
        if(strpos($redis->ping(),'PONG')!== false){
            return array(0,$redis);
        }
        return array(1,"");
    }

    function RedisQuit($red){
        $red->close();
    }
    
    function RedisSaveJson($red,$key,$json){
        $tmp=json_decode($json);
        // unset($tmp->{'status'});//unset: destroys the specified variables.
        $jsonstr = json_encode($tmp);
        if($red->set($key,$jsonstr)==1){
            $red->expire($key,$tmp->{'expired'});
            return 0;
        }else{
            return 1;
        }
    }

    function RedisGetJson($red,$key){
        if($red->exists($key)){
            $jsonstr = $red->get($key);
            return array(0,$jsonstr);
        }else{
            // list($ret,$jsonstr)=function($key);
            // if($ret!=0){
            //     return array($ret,"");
            // }
            // $ret=RedisSaveJson($red,$key,$jsonstr);
            // return array($ret,$jsonstr);
            return array(-1,"");
        }
    }
?>