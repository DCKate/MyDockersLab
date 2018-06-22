<?php
	// $ld_path=$_SERVER['LD_PATH']; 
    // putenv("LD_LIBRARY_PATH=$ld_path");
    
    function HttpHeaderParse(){
        if (!function_exists("getheaders"))
        {
            function getheaders()
            {
                $result = array();
                foreach($_SERVER as $key => $value)
                {
                    if (substr($key, 0, 5) == "HTTP_")
                    {
                        $key = str_replace(" ", "-", ucwords(strtolower(str_replace("_", " ", substr($key, 5)))));
                            $result[$key] = $value;
                    }
                }
                return $result;
            }
        }

        foreach (getheaders() as $headerName => $headerValue)
        {
            $a =array('key' => $headerName, 'value' => $headerValue);
            echo urldecode(json_encode($a));
            error_log($headerName." : ".$headerValue);
        }	
    }  

?>