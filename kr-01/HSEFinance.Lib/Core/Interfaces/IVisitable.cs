namespace HSEFinance.Lib.Core.Interfaces
{
    public interface IVisitable
    {
        void Accept(IVisitor visitor);
    }
}