using System;
using System.Collections.Generic;
using System.Linq;
using System.Runtime.InteropServices;
using System.Text;
using System.Threading.Tasks;

namespace RoBotUI.Bot.Utils
{
    public static class SPackHelper
    {
        public static Dictionary<string, object> Unpack(byte[] b, string template, params string[] names)
        {
            var split = template.Split(new []{ ' '},StringSplitOptions.RemoveEmptyEntries);

            var datas = new Dictionary<string,object>();

            int index = 0;
            for (int i  = 0; i < split.Length; i++)
            {
                var s = split[i];
                var parseChar = s.FirstOrDefault();

                var parseCount = 1;
                if (s.Length > 1)
                {
                    if (s[1] == '*')
                    {
                        /*if (datas.ContainsKey("len"))
                            parseCount = (ushort) datas["len"];
                        else*/
                            parseCount = b.Length - index;
                    }
                    else if (int.TryParse(s.Substring(1), out var parseNumber))
                        parseCount = parseNumber;
                }

                parseCount = Math.Min(parseCount, b.Length - index);

                var oldIndex = index;

                object data = null;

                switch (parseChar)
                {
                    case 'a':
                        data = new byte[parseCount];
                        Array.Copy(b,index,(byte[])data,0,parseCount);
                        index += parseCount;
                        break;
                    case 'Z':
                        data = Encoding.ASCII.GetString(b, index, parseCount).Trim(' ','\0');
                        index += parseCount;
                        break;
                    case 'x':
                        index += parseCount;
                        continue;
                    case 'v':
                        data = BitConverter.ToUInt16(b, index);
                        index += sizeof(UInt16);
                        break;
                    case 'V':
                        data = BitConverter.ToUInt32(b, index);
                        index += sizeof(UInt32);
                        break;
                    case 'l':
                        data = BitConverter.ToInt32(b, index);
                        index += sizeof(Int32);
                        break;
                    case 'c':
                        data = b[index];
                        index++;
                        break;
                    case 'C':
                        data = b[index];
                        index++;
                        break;
                }

                var name = names[i];
                datas[name] = data;

                //CLog.WriteLine("Old index {0}; Data type {1}; Data size {2}; New index {3}",oldIndex,parseChar,parseCount,index);
            }

            return datas;
        }
    }
}
