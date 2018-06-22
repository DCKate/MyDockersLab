<?php

class AES128{
    const CIPHER='AES-128-CBC';
    private $aeskey;
    private $aesiv;
    function __construct($key,$iv){
        $this->aeskey=$key;
        $this->aesiv=$iv;
    }
    public function aencrypt($data){
        $padding = 32 - (strlen($data) % 32);
        $data.= str_repeat(sprintf("%02d",$padding), $padding);
        $encrypt = openssl_encrypt($data, self::CIPHER, $this->aeskey, $options=OPENSSL_RAW_DATA, $this->aesiv);
        $encrypt_text = base64_encode($encrypt);
        return strtr($encrypt_text, '+/', '-_');
    }
    public function adecrypt($data)
    {
        $decrypt_text = base64_decode(strtr($data, '-_', '+/'));
        $decrypt = openssl_decrypt($decrypt_text, self::CIPHER, $this->aeskey, $options=OPENSSL_RAW_DATA, $this->aesiv);
        $trimnum = intval(substr($decrypt, -2));
        // return trim($decrypt,substr($decrypt, -1));
        return substr($decrypt, 0, 0-2*$trimnum);
    }

}

?>
