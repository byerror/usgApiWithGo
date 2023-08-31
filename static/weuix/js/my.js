$(function () {
    
    // var firstLoad = true;
    //创建MeScroll对象,内部已默认开启下拉刷新,自动执行up.callback,刷新列表数据;
    var mescroll = new MeScroll("body", { //id固定"body"
         //直接改mescroll.min.js里的源码，在callback: function(c)里location.reload
    });

  

    /* 这个在pc端没用
   $('.status-bar').doubleTap(() => {
                     
        $.toptip('doubleTap','success')
    });
    */








});