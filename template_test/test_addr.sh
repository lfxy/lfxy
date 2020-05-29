cat c.txt
if [ $? != 0 ]; then
exit 1
fi
sleep 1
addr="zk-0:2181"
echo "service_addr "$addr
