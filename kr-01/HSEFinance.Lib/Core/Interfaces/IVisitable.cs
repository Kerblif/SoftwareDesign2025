namespace HSEFinance.Lib.Core
{
    public interface IVisitable
    {
        void Accept(IVisitor visitor);
    }
}