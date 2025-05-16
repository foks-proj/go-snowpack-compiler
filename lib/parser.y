
%{
package lib

}

%union {
    file     *FileNode
    uniqueId *UniqueID
    uint     uint64
    rawval   string
}

%type <file> top
%type <uniqueId> fileID uniqueID


%%

top: 
    fileID
    {
        $$ = &FileNode{ Id: $1 }
    }
    ;

fileId:
    uniqueId SEMICOLON
    {
        $$ = $1
    }
    ;

uniqueId:
    TokenAt
    TokenUint64Val
    {
        $$ = &UniqueID{ Val: $2 }
    }
    ;

%%