#!/usr/bin/php
<?php
	require 'php_composer/vendor/autoload.php';
	require_once __DIR__ . '/pmodules/AES128.php';

    $Loader = new josegonzalez\Dotenv\Loader(__DIR__ . "/../conf/envrc");
    $Loader->parse();
    $Loader->toEnv();
    $akey=$_ENV['AES_CBC_KEY'];
    $aiv=$_ENV['AES_CBC_IV'];
    if(sizeof($akey)==0 || sizeof($aiv)==0){
        echo "error";
        return;
    }
    
    $uid=$argv[1];
    var_dump($argv);
    $encry=new AES128($akey,$aiv);
    $euid=$encry->aencrypt($uid);

    var_dump($euid);

?>
