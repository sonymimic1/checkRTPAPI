EXE="$1"
EXEPATH=$(dirname $0)
RUN="$EXEPATH/$EXE"
PIDOF=$(which pidof)
PIDFILE=pid

if [ -z $1 ]
then
        echo "do.sh {EXE} {options}"
        exit 1
fi

if [ ! -e $1 ]
then
        echo "$1 not exist!"
        exit 1
fi

start_exe()
{
        echo "Start $EXE : $RUN ..."
        nohup $RUN >> /dev/null 2> error.log &
        echo $! > $PIDFILE
        echo "... $! started successfully"
}

stop_exe()
{
        pidNUM=$(cat $PIDFILE)
        echo "Stop $EXE with pid $pidNUM..."
        if [ -z $pidNUM ]
        then
                echo "$EXE not found"
                return
        fi

        if kill $pidNUM
        then
                echo "$pidNUM is killed successfully"
        fi
        # rm $PIDFILE
}

ps_exe()
{
        pidNUM=$(cat $PIDFILE)
        echo "Status $EXE with pid $pidNUM"
        if [ -z $pidNUM ]
        then
                ps -ef | grep $EXE
        else
                ps -ef | grep $EXE | grep -f $PIDFILE
        fi
}

case $2 in
        "start")
                start_exe
                ;;
        "restart")
                stop_exe
                start_exe
                ;;
        "stop")
                stop_exe
                ;;
        "ps")
                ps_exe
                ;;
        *)
                echo "options: start / restart / stop / ps"
                ;;
esac


