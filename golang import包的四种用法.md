#### import包的四种用法

```bash
import "a/lib" // 通过包名libs调用导出接口 libs.Xxx()
import m "a/lib" // 通过别名m调用导出接口m.Xxx()
import . "a/lib" // .符号表示，对包lib中的导出接口的调用直接省略包名
import _ "a/lib" // _仅仅会导致lib执行初始化工作，如初始化全局变量，调用init函数。
```



