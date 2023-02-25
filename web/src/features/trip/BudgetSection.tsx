import React, {
  FC,
  useEffect,
  useState,
} from 'react';
import _flatten from "lodash/flatten";
import _get from "lodash/get";
import {
  ArrowLongRightIcon,
  ChatBubbleLeftEllipsisIcon,
  CurrencyDollarIcon,
  DocumentTextIcon,
  PencilIcon,
  PencilSquareIcon,
  PlusIcon,
  TrashIcon,
  XMarkIcon,
} from '@heroicons/react/24/solid';
import {
  EllipsisHorizontalCircleIcon
} from '@heroicons/react/24/outline';
import { Tooltip } from 'react-tooltip'

import Modal from '../../components/common/Modal';
import Dropdown from '../../components/common/Dropdown';

import {
  JSONPathBudgetAmount,
  getBudgetAmt,
  getBudgetItemPriceAmt,
  ContentIconOpts,
  getFlilghtPriceAmt,
  getItineraryContentPriceAmt,
  getLodgingPriceAmt,
  JSONPathPriceAmount,
  getContentColor,
  getTripContentForItineraryContent,
  BudgetItem,
  ItineraryList,
  ItineraryContent,
} from '../../lib/trips';
import {
  fmt,
  isEmptyDate,
  parseISO,
} from '../../lib/dates';
import {
  CommonCss,
  InputCss,
  TripBudgetCss,
  TripLogisticsCss
} from '../../assets/styles/global';
import {
  makeAddOp,
  makeRemoveOp,
  makeReplaceOp
} from '../../lib/jsonpatch';
import { Price } from '../../lib/common';





const ProgressBarId = "budget-progress-bar"

/////////////////////
// BudgetItemModal //
/////////////////////

interface BudgetItemModalProps {
  header?: string
  isOpen: boolean
  defaultTitle?: string
  defaultDesc?: string
  defaultAmount?: number
  onSubmit: (amount: Number|undefined, title: string, desc: string) => void
  onClose: () => void
}

const BudgetItemModal: FC<BudgetItemModalProps> = (props: BudgetItemModalProps) => {

  const [amount, setAmount] = useState<Number|undefined>(props.defaultAmount);
  const [title, setTitle] = useState(props.defaultTitle || "");
  const [desc, setDesc] = useState(props.defaultDesc || "");


  useEffect(() => {
    setAmount(props.defaultAmount)
    setTitle(props.defaultTitle || "")
    setDesc(props.defaultDesc || "")
  }, [props.defaultTitle, props.defaultDesc, props.defaultAmount])

  // Renderers
  const renderHeader = () => {
    return (
      <div className='flex justify-between items-center mb-4'>
        <div className='text-gray-800 font-bold text-xl'>
          {props.header}
        </div>
        <button type="button" onClick={() => {props.onClose()}}>
          <XMarkIcon className={CommonCss.Icon} />
        </button>
      </div>
    );
  }

  return (
    <Modal isOpen={props.isOpen}>
      <div className='p-5'>
        {renderHeader()}
        <div className={TripBudgetCss.PriceInputCtn}>
          <span className={TripBudgetCss.PriceInputLabel}>
            <CurrencyDollarIcon className={CommonCss.LeftIcon} />
            Amount
          </span>
          <input
            type="number"
            value={amount as any}
            onChange={(e) => setAmount(Number(e.target.value))}
            className={InputCss.Input}
          />
        </div>
        <div className={TripBudgetCss.PriceInputCtn}>
          <span className={TripBudgetCss.PriceInputLabel}>
            <ChatBubbleLeftEllipsisIcon className={CommonCss.LeftIcon} />
            Name
          </span>
          <input
            type="text"
            value={title}
            onChange={(e) => setTitle(e.target.value)}
            className={InputCss.Input}
          />
        </div>
        <div className={TripBudgetCss.PriceInputCtn}>
          <span className={TripBudgetCss.PriceInputLabel}>
            <DocumentTextIcon className={CommonCss.LeftIcon} />
            Description
          </span>
          <input
            type="text"
            value={desc}
            onChange={(e) => setDesc(e.target.value)}
            className={InputCss.Input}
          />
        </div>
        <div className='flex justify-around'>
          <button
            type='button'
            onClick={() => {
              props.onSubmit(amount, title, desc)
              props.onClose()
            }}
            className='py-2 px-4 bg-indigo-500 text-white font-bold rounded-full'
          >
            Submit
          </button>
        </div>
      </div>
    </Modal>
  );
}



/////////////////////
// EditBudgetModal //
/////////////////////


interface EditBudgetModalProps {
  budgetAmount: number
  isOpen: boolean
  onClose: () => void
  onSubmit: (newAmount: number) => void
}

const EditBudgetModal: FC<EditBudgetModalProps> = (props: EditBudgetModalProps) => {

  const [amount, setAmount] = useState(props.budgetAmount);

  // Renderers
  const renderHeader = () => {
    return (
      <div className='flex justify-between items-center mb-4'>
        <div className='text-gray-800 font-bold text-xl'>
          Set Budget
        </div>
        <button type="button" onClick={() => {props.onClose()}}>
          <XMarkIcon className={CommonCss.Icon} />
        </button>
      </div>
    );
  }

  return (
    <Modal isOpen={props.isOpen}>
      <div className='p-5'>
        {renderHeader()}
        <div className={TripBudgetCss.PriceInputCtn}>
          <span className={TripBudgetCss.PriceInputLabel}>
            <CurrencyDollarIcon className={CommonCss.LeftIcon} />
            Amount
          </span>
          <input
            type="number"
            value={amount as any}
            onChange={(e) => setAmount(Number(e.target.value))}
            className={InputCss.Input}
          />
        </div>
        <div className='flex justify-around'>
          <button
            type='button'
            onClick={() => {
              props.onSubmit(amount)
              props.onClose()
            }}
            className='py-2 px-4 bg-indigo-500 text-white font-bold rounded-full'
          >
            Submit
          </button>
        </div>
      </div>
    </Modal>
  );
}

///////////////////
// BudgetSection //
///////////////////

interface BudgetSectionProps {
  trip: any
  tripStateOnUpdate: any
}

const BudgetSection: FC<BudgetSectionProps> = (props: BudgetSectionProps) => {

  const [isAddBudgetModalOpen, setIsAddBudgetModalOpen] = useState(false);
  const [isEditBudgetItemModalOpen, setIsEditBudgetItemModalOpen] = useState(false);
  const [isEditBudgetModalOpen, setIsEditBudgetModalOpen] = useState(false);
  const [totalAmount, setTotalAmount] = useState(0);
  const [selectedBudgetItemIdx, setSelectedBudgetItemIdx] = useState<number|undefined>();
  const [selectedBudgetItem, setSelectedBudgetItem] = useState<BudgetItem|undefined>();

  const calculateTotalAmount = () => {
    let total = 0;

    Object.values(_get(props.trip, "flights", {}))
      .forEach((lod: any) => { total += getFlilghtPriceAmt(lod)});

    Object.values(_get(props.trip, "lodgings", {}))
      .forEach((lod: any) => { total += getLodgingPriceAmt(lod)});

    _get(props.trip, "itinerary", [])
      .forEach((l: ItineraryList) => {
        l.contents.forEach((ctnt: ItineraryContent) => {
          total += getItineraryContentPriceAmt(ctnt)
        })
      })

    _get(props.trip, "budget.items", [])
      .forEach((i: BudgetItem) => {
        total += getBudgetItemPriceAmt(i);
      })

    return total;
  }

  useEffect(() => {
    setTotalAmount(calculateTotalAmount())
  }, [props.trip])


  // Event Handlers

  const addNewBudgetItem = (amount: Number|undefined, title: string, desc: string) => {
    props.tripStateOnUpdate([
      makeAddOp("/budget/items/-", {
        title: title,
        desc: desc,
        price: { amount, currency: ""} as Price,
        labels: {},
        tags: {},
      } as BudgetItem)
    ]);
  }

  const deleteBudgetItem = (idx: number) => {
    props.tripStateOnUpdate([makeRemoveOp(`/budget/items/${idx}`, "")]);
  }

  const updateBudgetItem = (amount: Number|undefined, title: string, desc: string) => {
    props.tripStateOnUpdate([
      makeReplaceOp(`/budget/items/${selectedBudgetItemIdx}/title`, title),
      makeReplaceOp(`/budget/items/${selectedBudgetItemIdx}/desc`, desc),
      makeReplaceOp(`/budget/items/${selectedBudgetItemIdx}/${JSONPathPriceAmount}`, amount),
    ]);
  }

  const editBudgetItemOnClick = (bi: BudgetItem, idx: number) => {
    setSelectedBudgetItem(bi);
    setSelectedBudgetItemIdx(idx);
    setIsEditBudgetItemModalOpen(true);
  }

  const updateBudgetAmount = (amount: number) => {
    props.tripStateOnUpdate([
      makeReplaceOp(`/budget/${JSONPathBudgetAmount}`, amount),
    ]);
  }

  // Renderers

  const renderSummarySection = () => {
    // Progress Bar
    const renderProgressBar = () => {
      let pbstyle = { width: "100%" }
      const budgetAmount = getBudgetAmt(props.trip.budget);
      if (budgetAmount !== 0) {
        const width = Math.floor((totalAmount / budgetAmount) * 100);
        pbstyle.width = `${width}%`
      }
      return (
        <>
          <div className={TripBudgetCss.ProgressBarCtn} id={ProgressBarId}>
            <div className={TripBudgetCss.ProgressBarWrapper}>
              <div className={TripBudgetCss.ProgressBar} style={pbstyle} />
            </div>
          </div>
          <Tooltip
            anchorId={ProgressBarId}
            offset={1}
            place="bottom"
            content={`${totalAmount}/${budgetAmount}`}
          />
        </>
      );
    }

    return (
      <div className={TripBudgetCss.SummaryCtn}>
        <div className={TripBudgetCss.SpendingCtn}>
          <h2 className={TripBudgetCss.SpendingTitle}>Total Spending</h2>
          <h3 className={TripBudgetCss.SpendingAmount}>${totalAmount}</h3>
          {renderProgressBar()}
        </div>
        <div className={TripBudgetCss.OptsCtn}>
          <button type="button"
            onClick={() => setIsAddBudgetModalOpen(true)}
            className={TripBudgetCss.AddExpenseBtn}
          >
            <PlusIcon className={CommonCss.LeftIcon} />
            Add expense
          </button>
          <button
            type="button"
            className={TripBudgetCss.EditBudgetBtn}
            onClick={() => setIsEditBudgetModalOpen(true)}
          >
            <PencilIcon className={CommonCss.LeftIcon} />
            Edit Budget
          </button>
        </div>
      </div>
    )
  }

  const renderFlights = () => {
    return (
      <div className='mb-4'>
        <h4 className={TripBudgetCss.SubsectionTxt}>
          Flights
        </h4>
        {
          Object.values(_get(props.trip, "flights", {}))
          .map((flight: any) => {
            const Icon = ContentIconOpts["flight"];
            return (
              <div key={flight.id} className={TripBudgetCss.ItemCtn}>
                <div className={TripBudgetCss.ItemDescCtn}>
                  <span className={TripBudgetCss.FlightItemIcon}>
                    <Icon className={CommonCss.Icon} />
                  </span>
                  <div className="flex">
                    <p className={TripBudgetCss.ItemNameTxt}>
                      {flight.depart.departure.airport.code}
                    </p>
                    <ArrowLongRightIcon
                      className={TripLogisticsCss.FlightTransitLongArrow}
                    />
                    <p className={TripBudgetCss.ItemNameTxt}>
                      {flight.depart.arrival.airport.code}
                    </p>
                  </div>
                </div>
                <span className={TripBudgetCss.ItemPriceTxt}>
                  ${getFlilghtPriceAmt(flight)}
                </span>
              </div>
            );
          })
        }
      </div>
    );
  }

  const renderLodgings = () => {
    const dateFmt = "MMM dd";
    return (
      <div className='mb-4'>
        <h4 className={TripBudgetCss.SubsectionTxt}>
          Lodgings
        </h4>
        {
          Object.values(_get(props.trip, "lodgings", {}))
          .map((lod: any) => {
            const Icon = ContentIconOpts["hotel"];
            return (
              <div key={lod.id} className={TripBudgetCss.ItemCtn}>
                <div className={TripBudgetCss.ItemDescCtn}>
                  <span className={TripBudgetCss.LodgingItemIcon}>
                    <Icon className={CommonCss.Icon} />
                  </span>
                  <div>
                    <p className={TripBudgetCss.ItemNameTxt}>
                      {lod.place.name}
                    </p>
                    <p className={TripBudgetCss.LodgingDatesTxt}>
                      {isEmptyDate(lod.checkinTime) ? null
                        : fmt(parseISO(lod.checkinTime as string), dateFmt)}
                      {isEmptyDate(lod.checkoutTime) ? null :
                        " - " + fmt(parseISO(lod.checkoutTime as string), dateFmt)}
                    </p>
                  </div>
                </div>
                <span className={TripBudgetCss.ItemPriceTxt}>
                  ${getLodgingPriceAmt(lod)}
                </span>
              </div>
            );
          })
        }
      </div>
    );
  }

  const renderItinerary = () => {
    const itinerary = _flatten(_get(props.trip, "itinerary", [])
      .map((l: ItineraryList) => {
        return l.contents.map((itinCtnt: ItineraryContent, idx: number) => {
          const amt = getItineraryContentPriceAmt(itinCtnt);
          if (amt === 0) {
            return null;
          }
          const ctnt = getTripContentForItineraryContent(
            props.trip, itinCtnt.tripContentListId, itinCtnt.tripContentId
          );
          const color = getContentColor(l);
          return (
            <div key={itinCtnt.id} className={TripBudgetCss.ItemCtn}>
              <div className={TripBudgetCss.ItemDescCtn}>
                <div
                  className={TripBudgetCss.ItinItemIcon}
                  style={{backgroundColor: color}}
                >
                  {idx + 1}
                </div>
                <div>
                  <p className={TripBudgetCss.ItemNameTxt}>{ctnt.title}</p>
                  <p className={TripBudgetCss.ItemDescTxt}>
                    {fmt(parseISO(l.date as string), "eee, MM/dd")}
                  </p>
                </div>
              </div>
              <span className={TripBudgetCss.ItemPriceTxt}>
                ${amt}
              </span>
            </div>
          );
        })
      })
    )
    .filter((item: any) => item !== null);
    return (
      <div className='mb-4'>
        <h4 className={TripBudgetCss.SubsectionTxt}>
          Itinerary
        </h4>
        {itinerary as any}
      </div>
    );
  }

  const renderBudgetItemSettingsDropdown = (bi: BudgetItem, idx: number) => {
    const opts = [
      <button
        type='button'
        className={CommonCss.DropdownBtn}
        onClick={() => { editBudgetItemOnClick(bi, idx) }}
      >
        <PencilSquareIcon className={CommonCss.LeftIcon} />
        Update
      </button>,
      <button
        type='button'
        className={CommonCss.DeleteBtn}
        onClick={() => deleteBudgetItem(idx)}
      >
        <TrashIcon className={CommonCss.LeftIcon} />
        Delete
      </button>,
    ];
    const menu = (<EllipsisHorizontalCircleIcon className={CommonCss.DropdownIcon} />);
    return <Dropdown menu={menu} opts={opts} />
  }

  const renderBudgetItems = () => {
    return (
      <div>
        <h4 className={TripBudgetCss.SubsectionTxt}>
          Custom
        </h4>
        {
          _get(props.trip, "budget.items", [])
          .map((bi: BudgetItem, idx: number) => {
            const amt = getBudgetItemPriceAmt(bi);
            return (
              <div key={bi.id}
                className={TripBudgetCss.ItemCtn}
              >
                <div className={TripBudgetCss.ItemDescCtn}>
                  <div className={TripBudgetCss.BudgetItemIcon}>
                    {idx + 1}
                  </div>
                  <div>
                    <p className={TripBudgetCss.ItemNameTxt}>{bi.title}</p>
                    <p className={TripBudgetCss.ItemDescTxt}>{bi.desc}</p>
                  </div>
                </div>
                <div className='flex items-center'>
                  <span className={TripBudgetCss.ItemPriceTxt + " mr-2"}>
                    ${amt}
                  </span>
                  {renderBudgetItemSettingsDropdown(bi, idx)}
                </div>
              </div>
            );
          })
        }
      </div>
    );
  }

  return (
    <div className='p-5'>
      {renderSummarySection()}
      {renderFlights()}
      {renderLodgings()}
      {renderItinerary()}
      {renderBudgetItems()}
      <BudgetItemModal
        header='Add custom expense'
        isOpen={isAddBudgetModalOpen}
        onSubmit={addNewBudgetItem}
        onClose={() => {setIsAddBudgetModalOpen(false)}}
      />
      <BudgetItemModal
        header='Edit expense'
        isOpen={isEditBudgetItemModalOpen}

        defaultTitle={_get(selectedBudgetItem, "title")}
        defaultDesc={_get(selectedBudgetItem, "desc")}
        onSubmit={updateBudgetItem}
        onClose={() => {setIsEditBudgetItemModalOpen(false)}}
      />
      <EditBudgetModal
        budgetAmount={getBudgetAmt(props.trip.budget)}
        isOpen={isEditBudgetModalOpen}
        onClose={() => setIsEditBudgetModalOpen(false)}
        onSubmit={updateBudgetAmount}
      />
    </div>
  );
}


export default BudgetSection;
